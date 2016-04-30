package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path/filepath"
	"github.com/TakatoshiMaeda/kinu/engine"
	"errors"
	"encoding/json"
	"github.com/TakatoshiMaeda/kinu/logger"
	"net/http/pprof"
)

var ErrInvalidImageExt = errors.New("supported image type is only jpg/jpeg")

type ErrInvalidRequest struct {
	error
	Message string
}

func (e *ErrInvalidRequest) Error() string { return e.Message }

func main() {
	logger.Debug("Kinu Booting...")
	engine.Initialize()
	defer engine.Finalize()

	router := httprouter.New()
	router.GET("/:type/:id/:geometry/:filename", GetImageHandler)

	router.POST("/upload", UploadImageHandler)
	router.POST("/sandbox", UploadImageToSandboxHandler)
	router.POST("/sandbox/apply", ApplyFromSandboxHandler)

	router.GET("/debug/pprof/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { pprof.Index(w, r) })
	router.GET("/debug/pprof/cmdline", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { pprof.Cmdline(w, r) })
	router.GET("/debug/pprof/profile", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { pprof.Profile(w, r) })
	router.GET("/debug/pprof/symbol", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { pprof.Symbol(w, r) })

	logger.Debug("Started Kinu.")
	logger.Fatal(http.ListenAndServe(":8080", router))
}

func SetContentType(w http.ResponseWriter, filename string) error {
	ext := ExtractExtension(filepath.Ext(filename))
	if IsValidImageExt(ext) {
		w.Header().Set("Content-Type", "image/" + ext)
		return nil
	} else {
		return ErrInvalidImageExt
	}
}

func RespondBadRequest(w http.ResponseWriter, reason string) {
	w.Header().Set("X-BadRequest-Reason", reason)
	w.WriteHeader(http.StatusBadRequest)
}

func RespondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func RespondInternalServerError(w http.ResponseWriter, err error) {
	logger.ErrorDebug(err)
	w.WriteHeader(http.StatusInternalServerError)
}

type UploadResult struct {
	ImageType string
	ImageId   string
}

func RespondImageUploadSuccessJson(w http.ResponseWriter, imageType string, imageId string) {
	json, err := json.Marshal(&UploadResult{ImageType: imageType, ImageId: imageId})
	if err != nil {
		RespondInternalServerError(w, err)
	}
	RespondJson(w, json)
}

func RespondJson(w http.ResponseWriter, json []byte) {
	w.Write(json)
}

func RespondImage(w http.ResponseWriter, bytes []byte) {
	w.Write(bytes)
}