package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/TakatoshiMaeda/kinu/storage"
	"sync"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/Sirupsen/logrus"
	"strconv"
)

const SANDBOX_IMAGE_TYPE = "__sandbox__"

type ErrAttachFromSandbox struct {
	error
	Errors []error
}

func (e *ErrAttachFromSandbox) Error() string {
	messages := "Image attach from sandbox error. cause, "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i + 1) + ". " + err.Error() + "  "
	}
	return messages
}

func UploadImageToSandboxHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(0)

	imageId := uuid.NewV4().String()

	file, _, err := r.FormFile("image")
	if err != nil {
		RespondBadRequest(w, "invalid file")
		return
	}

	err = UploadImage(SANDBOX_IMAGE_TYPE, imageId, file)
	if err != nil {
		if _, ok := err.(*ErrInvalidRequest); ok {
			RespondBadRequest(w, err.Error())
		} else {
			RespondInternalServerError(w, err)
		}
	}

	RespondImageUploadSuccessJson(w, SANDBOX_IMAGE_TYPE, imageId)

	logger.WithFields(logrus.Fields{
		"path": r.URL.Path,
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

	sandboxImageMetadata := NewImageMetadata(SANDBOX_IMAGE_TYPE, sandboxId)

	st, err := storage.Open()
	if err != nil {
		RespondInternalServerError(w, err)
		return
	}

	items, err := st.List(sandboxImageMetadata.BasePath())
	if err != nil {
		RespondInternalServerError(w, err)
		return
	}

	applyImageMetadata := NewImageMetadata(imageType, imageId)

	wg := sync.WaitGroup{}
	errs := make(chan error, len(items))
	for _, item := range items {
		wg.Add(1)
		go func(item storage.StorageItem){
			defer wg.Done()
			st, err := storage.Open()
			if err != nil {
				errs <- logger.ErrorDebug(err)
				return
			}

			err = st.Move(item.Key(), applyImageMetadata.FilePath(item.ImageSize(), item.Extension()))
			if err != nil {
				errs <- logger.ErrorDebug(err)
				return
			}

			errs <- nil
		}(item)
	}
	wg.Wait()

	close(errs)

	errors := make([]error, 0)
	for err := range errs {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		RespondInternalServerError(w, &ErrAttachFromSandbox{Errors: errors})
		return
	}

	RespondImageUploadSuccessJson(w, imageType, imageId)

	logger.WithFields(logrus.Fields{
		"path": r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
	}).Info("success")
}
