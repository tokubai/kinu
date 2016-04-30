package main

import (

	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/Sirupsen/logrus"
	"errors"
	"time"
)

var ErrImageNotFound = errors.New("image not found.")

func GetImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger.WithFields(logrus.Fields{
		"path": r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("start processing.")

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
		if err == ErrImageNotFound {
			RespondNotFound(w)
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}
	logger.TrackResult("fetch image from storage", imageFetchStartTime)

	resizeStartTime := time.Now()
	if imageGetRequest.NeedsOriginalImage() {
		RespondImage(w, originalImage)
		return
	}

	resizedImage, err := resizer.Run(originalImage, imageGetRequest.ToResizeOption())
	if err != nil {
		RespondInternalServerError(w, err)
		return
	}
	logger.TrackResult("resize image", resizeStartTime)

	RespondImage(w, resizedImage)

	logger.WithFields(logrus.Fields{
		"path": r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("process success.")
}
