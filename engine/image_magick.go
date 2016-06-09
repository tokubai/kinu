package engine

import (
	"fmt"
	"github.com/TakatoshiMaeda/kinu/logger"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type ImageMagickEngine struct {
	ResizeEngine

	mw                *imagick.MagickWand
	opened            bool
	originalImageBlob []byte

	heightSizeHint, widthSizeHint int
}

func newImageMagickEngine(image []byte) (e *ImageMagickEngine) {
	return &ImageMagickEngine{originalImageBlob: image}
}

func (e *ImageMagickEngine) SetSizeHint(width int, height int) {
	e.heightSizeHint = height
	e.widthSizeHint = width
}

func (e *ImageMagickEngine) Open() error {
	e.mw = imagick.NewMagickWand()
	if e.heightSizeHint > 0 && e.widthSizeHint > 0 {
		e.mw.SetOption("jpeg:size", fmt.Sprintf("%dx%d", e.heightSizeHint, e.widthSizeHint))
	}
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
	orientation := e.mw.GetImageOrientation()
	if orientation != imagick.ORIENTATION_UNDEFINED && orientation != imagick.ORIENTATION_TOP_LEFT {
		ok := e.mw.AutoOrientImage()
		if ok != nil {
			return nil, ok
		}
	}
	ok := e.mw.StripImage()
	if ok != nil {
		return nil, ok
	}
	return e.mw.GetImageBlob(), nil
}
