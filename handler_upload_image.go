package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(0)

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

	file, _, err := r.FormFile("image")
	if err != nil {
		RespondBadRequest(w, "invalid file")
		return
	}

	err = UploadImage(imageType, imageId, file)
	if err != nil {
		if _, ok := err.(*ErrInvalidRequest); ok {
			RespondBadRequest(w, err.Error())
		} else {
			RespondInternalServerError(w, err)
		}
	}

	RespondImageUploadSuccessJson(w, imageType, imageId)
}
