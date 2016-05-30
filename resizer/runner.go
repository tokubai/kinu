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
	image      []byte
	option     *ResizeOption
	resultChan chan *ResizeResult
}

type ResizeResult struct {
	image []byte
	err   error
}

var (
	ResizeWorkerRunningLimitMaxNum int
	ResizeWorkerWaitPoolMaxNum     int

	ErrTooManyRunningResizeWorker = errors.New("Too many running resize worker error.")

	resizeWorkerWaitLimiter chan bool
	resizeWorkerLimiter     chan bool

	resizeRequestDispatcher = make(chan *ResizeRequest)
)

const (
	DEFAULT_QUALITY = 70
)

func Run(image []byte, option *ResizeOption) (resizedImage []byte, err error) {
	if IsFreeResizeWorkerAvailable() {
		result := <-dispatch(image, option)
		return result.image, result.err
	} else {
		return nil, ErrTooManyRunningResizeWorker
	}
}

func IsFreeResizeWorkerAvailable() bool {
	return len(resizeWorkerWaitLimiter) < ResizeWorkerWaitPoolMaxNum
}

func init() {
	maxNum := os.Getenv("KINU_RESIZE_WORKER_MAX_SIZE")
	if len(maxNum) != 0 {
		num, err := strconv.Atoi(maxNum)
		if err != nil {
			panic(err)
		}
		ResizeWorkerRunningLimitMaxNum = num
	} else {
		ResizeWorkerRunningLimitMaxNum = runtime.NumCPU() * 15
	}
	resizeWorkerLimiter = make(chan bool, ResizeWorkerRunningLimitMaxNum)

	waitPool := os.Getenv("KINU_RESIZE_WAIT_POOL_SIZE")
	if len(waitPool) != 0 {
		num, err := strconv.Atoi(waitPool)
		if err != nil {
			panic(err)
		}
		ResizeWorkerWaitPoolMaxNum = num
	} else {
		ResizeWorkerWaitPoolMaxNum = runtime.NumCPU() * 20
	}
	resizeWorkerWaitLimiter = make(chan bool, ResizeWorkerWaitPoolMaxNum)

	runWorker()
}

func runWorker() {
	go func() {
		for r := range resizeRequestDispatcher {
			go func() {
				resizeWorkerWaitLimiter <- true
				defer func() { <-resizeWorkerWaitLimiter }()
				work(r.image, r.option, r.resultChan)
			}()
		}
	}()
}

func work(image []byte, option *ResizeOption, resultChan chan *ResizeResult) {
	resizeWorkerLimiter <- true
	defer func() {
		<-resizeWorkerLimiter
	}()

	logger.WithFields(logrus.Fields{
		"resize_worker_num":          len(resizeWorkerLimiter),
		"resize_worker_max_num":      ResizeWorkerRunningLimitMaxNum,
		"resize_worker_pool_num":     len(resizeWorkerWaitLimiter),
		"resize_worker_pool_max_num": ResizeWorkerWaitPoolMaxNum,
	}).Debug("resize worker status")

	resultChan <- Resize(image, option)
}

func dispatch(image []byte, option *ResizeOption) (resultChan chan *ResizeResult) {
	request := &ResizeRequest{
		image:      image,
		option:     option,
		resultChan: make(chan *ResizeResult, 1),
	}
	resizeRequestDispatcher <- request
	return request.resultChan
}

type ResizeOption struct {
	Width         int
	Height        int
	NeedsAutoCrop bool
	Quality       int

	SizeHintWidth  int
	SizeHintHeight int
}

func (o *ResizeOption) HasSizeHint() bool {
	return o.SizeHintHeight > 0 && o.SizeHintWidth > 0
}
