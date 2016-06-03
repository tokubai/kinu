package main

import (
	"encoding/json"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type WorkerStats struct {
	WorkerTotalNum          int    `json:"total_worker_num"`
	ResizeRequestBufferSize int    `json:"resize_request_buffer_size"`
	StackedResizeRequestNum int    `json:"stacked_resize_request_num"`
	TooManyResizeRequest    bool   `json:"too_many_resize_request"`
}

func WorkerStatsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if !resizer.IsWorkerMode {
		RespondBadRequest(w, "current kinu mode is not worker mode")
		return
	}

	js, err := json.Marshal(&WorkerStats{
		WorkerTotalNum:          resizer.ResizeWorkerSize,
		ResizeRequestBufferSize: resizer.ResizeWorkerWaitBufferNum,
		StackedResizeRequestNum: resizer.RequestPayloadLen(),
		TooManyResizeRequest:    !resizer.CanResizeRequest(),
	})

	if err != nil {
		RespondInternalServerError(w, err)
	}

	RespondJson(w, js)
}
