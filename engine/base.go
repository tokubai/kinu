package engine

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tokubai/kinu/logger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type ResizeEngine interface {
	Open() error
	Close()

	SetSizeHint(width, height int)
	SetFormat(format string)
	SetCompressionQuality(quality int)

	GetImageHeight() int
	GetImageWidth() int

	RemoveAlpha() error
	Resize(width int, height int) error
	Crop(width int, height int, startX int, startY int) error
	Generate() ([]byte, error)
}

var (
	AvailableEngines       = []string{"ImageMagick"}
	ErrUnknownResizeEngine = errors.New("specify unknown resize engine.")
	selectedEngineType     string
)

func init() {
	selectedEngineType = os.Getenv("KINU_RESIZE_ENGINE")
	if len(selectedEngineType) == 0 {
		panic("must specify KINU_RESIZE_ENGINE system environment.")
	}

	var isAvailableEngine bool
	for _, engineType := range AvailableEngines {
		if selectedEngineType == engineType {
			isAvailableEngine = true
		}
	}

	if !isAvailableEngine {
		panic("unknown KINU_RESIZE_ENGINE " + selectedEngineType + ".")
	}

	logger.WithFields(logrus.Fields{
		"resize_engine_type": selectedEngineType,
	}).Info("setup resize engine")
}

func New(image []byte) (ResizeEngine, error) {
	switch selectedEngineType {
	case "ImageMagick":
		return newImageMagickEngine(image), nil
	default:
		return nil, ErrUnknownResizeEngine
	}
}

func Initialize() {
	if selectedEngineType == "ImageMagick" {
		imagick.Initialize()
	}
}

func Finalize() {
	if selectedEngineType == "ImageMagick" {
		imagick.Terminate()
	}
}
