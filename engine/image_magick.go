package engine
import (
	"gopkg.in/gographics/imagick.v2/imagick"
	"github.com/TakatoshiMaeda/kinu/logger"
)

type ImageMagickEngine struct {
	ResizeEngine

	mw *imagick.MagickWand
	opened bool
	originalImageBlob []byte
}

func newImageMagickEngine(image []byte) (e *ImageMagickEngine) {
	return &ImageMagickEngine{originalImageBlob: image}
}

func (e *ImageMagickEngine) Open() error {
	e.mw = imagick.NewMagickWand()
	err := e.mw.ReadImageBlob(e.originalImageBlob)
	if err != nil {
		return logger.ErrorDebug(err)
	} else {
		e.opened = true
	}
	return nil
}

func (e *ImageMagickEngine) Close() {
	if e.opened {
		e.mw.Destroy()
	}
}

func (e *ImageMagickEngine) GetImageHeight() int {
	return int(e.mw.GetImageHeight())
}

func (e *ImageMagickEngine) GetImageWidth() int {
	return int(e.mw.GetImageWidth())
}

func (e *ImageMagickEngine) Resize(width int, height int) error {
	return e.mw.ResizeImage(uint(width), uint(height), imagick.FILTER_LANCZOS, 1.0)
}

func (e *ImageMagickEngine) Crop(width int, height int, startX int, startY int) error {
	return e.mw.CropImage(uint(width), uint(height), startX, startY)
}

func (e *ImageMagickEngine) Generate() ([]byte, error) {
	return e.mw.GetImageBlob(), nil
}
