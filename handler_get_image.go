package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/TakatoshiMaeda/kinu/storage"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func GetImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := SetContentType(w, ps.ByName("filename"))
	if err != nil {
		if err == ErrInvalidImageExt {
			RespondBadRequest(w, err.Error())
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}

	imageGetRequest, err := NewImageGetRequest(ps)
	if err != nil {
		if _, ok := err.(*ErrInvalidRequest); ok {
			RespondBadRequest(w, err.Error())
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}

	imageFetchStartTime := time.Now()
	originalImage, err := imageGetRequest.FetchImage()
	if err != nil {
		if err == storage.ErrImageNotFound {
			RespondNotFound(w)
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}
	logger.TrackResult("fetch image from storage", imageFetchStartTime)

	if imageGetRequest.NeedsOriginalImage() {
		RespondImage(w, originalImage)
		return
	}

	resizeStartTime := time.Now()
	resizedImage, err := resizer.Run(originalImage, imageGetRequest.ToResizeOption())
	if err != nil {
		if err == resizer.ErrTooManyRunningResizeWorker {
			RespondServiceUnavailable(w, err)
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}
	logger.TrackResult("resize image", resizeStartTime)

	RespondImage(w, resizedImage)

	logger.WithFields(logrus.Fields{
		"path":   r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("success")
}
