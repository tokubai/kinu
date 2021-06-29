package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/tokubai/kinu/logger"
	"github.com/tokubai/kinu/resource"
)

const SANDBOX_IMAGE_TYPE = "__sandbox__"

func UploadImageToSandboxHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(0)

	uid := uuid.NewV4()
	imageId := uid.String()

	file, _, err := r.FormFile("image")
	if err != nil {
		RespondBadRequest(w, "invalid file")
		return
	}

	err = resource.New(SANDBOX_IMAGE_TYPE, imageId).Store(file)
	if err != nil {
		if _, ok := err.(*ErrInvalidRequest); ok {
			RespondBadRequest(w, err.Error())
		} else {
			RespondInternalServerError(w, err)
		}
	}

	RespondImageUploadSuccessJson(w, SANDBOX_IMAGE_TYPE, imageId)

	logger.WithFields(logrus.Fields{
		"path":   r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("success")
}

func ApplyFromSandboxHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(0)

	sandboxId := r.FormValue("sandbox_id")
	if len(sandboxId) == 0 {
		RespondBadRequest(w, "required sandbox_id")
		return
	}

	imageType := r.FormValue("name")
	if len(imageType) == 0 {
		RespondBadRequest(w, "required name parameter")
		return
	}

	imageId := r.FormValue("id")
	if len(imageId) == 0 {
		RespondBadRequest(w, "required id parameter")
		return
	}

	err := resource.New(SANDBOX_IMAGE_TYPE, sandboxId).MoveTo(imageType, imageId)

	if err != nil {
		RespondInternalServerError(w, err)
		return
	}

	RespondImageUploadSuccessJson(w, imageType, imageId)

	logger.WithFields(logrus.Fields{
		"path":   r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("success")
}
