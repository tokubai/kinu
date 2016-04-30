package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/TakatoshiMaeda/kinu/storage"
	"sync"
	"github.com/TakatoshiMaeda/kinu/logger"
)

const SANDBOX_IMAGE_TYPE = "__sandbox__"

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
}

func ApplyFromSandboxHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")

	r.ParseMultipartForm(0)

	sandboxId := r.FormValue("sandbox_id")
	if len(sandboxId) == 0 {
		RespondBadRequest(w, "required sandbox_id")
		return
	}

	imageType := r.FormValue("model")
	if len(imageType) == 0 {
		RespondBadRequest(w, "required model parameter")
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

	// FIXME: Error Handling
	applyImageMetadata := NewImageMetadata(imageType, imageId)
	wg := sync.WaitGroup{}
	for _, item := range items {
		wg.Add(1)
		go func(item storage.StorageItem){
			defer wg.Done()
			st, err := storage.Open()
			if err != nil {
				logger.ErrorDebug(err)
			}

			err = st.Move(item.Key(), applyImageMetadata.FilePath(item.ImageSize(), item.Extension()))
			if err != nil {
				logger.ErrorDebug(err)
			}
		}(item)
	}
	wg.Wait()

	RespondImageUploadSuccessJson(w, imageType, imageId)
}
