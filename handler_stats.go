package main

import (
	"encoding/json"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type WorkerStats struct {
	Version                 string `json:"version"`
	WorkerTotalNum          int    `json:"total_worker_num"`
	ResizeRequestBufferSize int    `json:"resize_request_buffer_size"`
	StackedResizeRequestNum int    `json:"stacked_resize_request_num"`
	TooManyResizeRequest    bool   `json:"too_many_resize_request"`
}

type Stats struct {
	Version string `json:"version"`
}

func StatsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var js []byte
	var err error

	if resizer.IsWorkerMode {
		js, err = json.Marshal(&WorkerStats{
			Version:                 VERSION,
			WorkerTotalNum:          resizer.ResizeWorkerSize,
			ResizeRequestBufferSize: resizer.ResizeWorkerWaitBufferNum,
			StackedResizeRequestNum: resizer.RequestPayloadLen(),
			TooManyResizeRequest:    !resizer.CanResizeRequest(),
		})
	} else {
		js, err = json.Marshal(&Stats{
			Version: VERSION,
		})
	}

	if err != nil {
		RespondInternalServerError(w, err)
	}

	RespondJson(w, js)
}
