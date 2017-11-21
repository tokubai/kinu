package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/tokubai/kinu/logger"
	"github.com/tokubai/kinu/resizer"
	"github.com/tokubai/kinu/resource"
	"github.com/tokubai/kinu/storage"
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

	request, err := NewImageGetRequest(ps)
	if err != nil {
		if _, ok := err.(*resizer.ErrInvalidGeometry); ok {
			RespondBadRequest(w, err.Error())
		} else if _, ok := err.(*resizer.ErrInvalidGeometryOrderRequest); ok {
			RespondNotFound(w)
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}

	resource := resource.New(request.Category, request.Id)

	imageFetchStartTime := time.Now()
	image, err := resource.Fetch(request.Geometry)
	if err != nil {
		if err == storage.ErrImageNotFound {
			RespondNotFound(w)
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}
	logger.TrackResult("fetch image from storage", imageFetchStartTime)

	if request.Geometry.NeedsOriginalImage {
		RespondImage(w, image.Body)
		return
	}

	if len(request.Geometry.MiddleImageSize) != 0 {
		RespondImage(w, image.Body)
		return
	}

	resizeStartTime := time.Now()
	resizeOption := request.Geometry.ToResizeOption()
	resizeOption.SizeHintHeight = image.Height
	resizeOption.SizeHintWidth = image.Width
	resizeOption.Format = request.Extension
	resizedImage, err := resizer.Run(image.Body, resizeOption)
	if err != nil {
		if err == resizer.ErrTooManyRunningResizeWorker {
			RespondServiceUnavailable(w, err)
		} else {
			RespondInternalServerError(w, err)
		}
		return
	}
	logger.TrackResult("resize image", resizeStartTime)

	if request.Extension == "data" {
		RespondDataURI(w, resizedImage)
	} else {
		RespondImage(w, resizedImage)
	}

	logger.WithFields(logrus.Fields{
		"path":   r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("success")
}
