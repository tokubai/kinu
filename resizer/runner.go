package resizer

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/TakatoshiMaeda/kinu/logger"
	"os"
	"runtime"
	"strconv"
)

type ResizeRequest struct {
	image         []byte
	option        *ResizeOption
	resultPayload chan *ResizeResult
}

type ResizeResult struct {
	image []byte
	err   error
}

var (
	IsWorkerMode bool

	ResizeWorkerSize          int
	ResizeWorkerWaitBufferNum int

	ErrTooManyRunningResizeWorker = errors.New("Too many running resize worker error.")

	requestPayload chan *ResizeRequest
)

const (
	DEFAULT_QUALITY = 70
)

func Run(image []byte, option *ResizeOption) (resizedImage []byte, err error) {
	if !IsWorkerMode {
		result := Resize(image, option)
		return result.image, result.err
	}

	if CanResizeRequest() {
		request := &ResizeRequest{image: image, option: option, resultPayload: make(chan *ResizeResult, 1)}
		requestPayload <- request
		result := <-request.resultPayload
		return result.image, result.err
	} else {
		return nil, ErrTooManyRunningResizeWorker
	}
}

func RequestPayloadLen() int {
	return len(requestPayload)
}

func CanResizeRequest() bool {
	return len(requestPayload) < ResizeWorkerWaitBufferNum
}

func init() {
	IsWorkerMode = (len(os.Getenv("KINU_RESIZE_WORKER_MODE")) != 0)
	if !IsWorkerMode {
		return
	}

	maxNum := os.Getenv("KINU_RESIZE_WORKER_MAX_SIZE")
	if len(maxNum) != 0 {
		num, err := strconv.Atoi(maxNum)
		if err != nil {
			panic(err)
		}
		ResizeWorkerSize = num
	} else {
		ResizeWorkerSize = runtime.NumCPU() * 10
	}

	waitPool := os.Getenv("KINU_RESIZE_WORKER_WAIT_BUFFER")
	if len(waitPool) != 0 {
		num, err := strconv.Atoi(waitPool)
		if err != nil {
			panic(err)
		}
		ResizeWorkerWaitBufferNum = num
	} else {
		ResizeWorkerWaitBufferNum = ResizeWorkerSize * 3
	}

	logger.WithFields(logrus.Fields{
		"worker_size":        ResizeWorkerSize,
		"resize_wait_buffer": ResizeWorkerWaitBufferNum,
	}).Info("set worker config")

	requestPayload = make(chan *ResizeRequest, ResizeWorkerWaitBufferNum)

	runWorker()
}

func runWorker() {
	for i := 1; i <= ResizeWorkerSize; i++ {
		go worker(i, requestPayload)
	}

	logger.WithFields(logrus.Fields{
		"worker_size": ResizeWorkerSize,
	}).Info("resize worker started")
}

func worker(id int, requests <-chan *ResizeRequest) {
	logger.WithFields(logrus.Fields{
		"worker_id": id,
	}).Debug("launch resize worker")

	for r := range requests {
		logger.WithFields(logrus.Fields{
			"worker_id": id,
		}).Debug("processing resize from worker")

		r.resultPayload <- Resize(r.image, r.option)
	}
}

type ResizeOption struct {
	Width            int
	Height           int
	NeedsAutoCrop    bool
	NeedsManualCrop  bool
	CropWidthOffset  int
	CropHeightOffset int
	CropWidth        int
	CropHeight       int
	AssumptionWidth  int
	Quality          int

	SizeHintWidth  int
	SizeHintHeight int
}

func (o *ResizeOption) HasSizeHint() bool {
	return o.SizeHintHeight > 0 && o.SizeHintWidth > 0
}
