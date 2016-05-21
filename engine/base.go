package engine

import (
	"errors"
	"gopkg.in/gographics/imagick.v2/imagick"
	"os"
)

type ResizeEngine interface {
	Open() error
	Close()

	SetResizeSize(width, height int)

	GetImageHeight() int
	GetImageWidth() int

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
