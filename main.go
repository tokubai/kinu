package main

import (
	"encoding/json"
	"errors"
	"github.com/TakatoshiMaeda/kinu/engine"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrInvalidImageExt = errors.New("supported image type is only jpg/jpeg")
	listenPort         = "80"
)

type ErrInvalidRequest struct {
	error
	Message string
}

type ErrInvalidGeometryOrderRequest struct {
	error
	Message string
}

func (e *ErrInvalidRequest) Error() string { return e.Message }

func init() {
	if len(os.Getenv("KINU_PORT")) != 0 {
		listenPort = os.Getenv("KINU_PORT")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	engine.Initialize()
	defer engine.Finalize()

	router := httprouter.New()

	if os.Getenv("KINU_DEBUG") == "1" {
		router.GET("/debug/pprof/*pprof", HandlePprof)
	}

	router.GET("/images/:type/:geometry/:filename", GetImageHandler)

	router.POST("/upload", UploadImageHandler)
	router.POST("/sandbox", UploadImageToSandboxHandler)
	router.POST("/sandbox/attach", ApplyFromSandboxHandler)

	logger.Fatal(http.ListenAndServe(":"+listenPort, router))
}

func HandlePprof(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch p.ByName("pprof") {
	case "/cmdline":
		pprof.Cmdline(w, r)
	case "/profile":
		pprof.Profile(w, r)
	case "/symbol":
		pprof.Symbol(w, r)
	default:
		pprof.Index(w, r)
	}
}

func SetContentType(w http.ResponseWriter, filename string) error {
	ext := ExtractExtension(filepath.Ext(filename))
	if IsValidImageExt(ext) {
		w.Header().Set("Content-Type", "image/"+ext)
		return nil
	} else {
		return ErrInvalidImageExt
	}
}

func RespondBadRequest(w http.ResponseWriter, reason string) {
	w.Header().Set("X-Kinu-BadRequest-Reason", reason)
	w.WriteHeader(http.StatusBadRequest)
}

func RespondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func RespondInternalServerError(w http.ResponseWriter, err error) {
	logger.ErrorDebug(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func RespondServiceUnavailable(w http.ResponseWriter, err error) {
	logger.ErrorDebug(err)
	w.WriteHeader(http.StatusServiceUnavailable)
}

type UploadResult struct {
	ImageType string `json:"name"`
	ImageId   string `json:"id"`
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
