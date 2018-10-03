package engine

import (
	"fmt"

	"github.com/tokubai/kinu/logger"
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

func (e *ImageMagickEngine) SetFormat(format string) {
	if format == "data" {
		e.mw.SetImageFormat("jpeg")
	} else {
		e.mw.SetImageFormat(format)
	}
}

func (e *ImageMagickEngine) SetCompressionQuality(quality int) {
	e.mw.SetImageCompressionQuality(uint(quality))
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

func (e *ImageMagickEngine) RemoveAlpha() error {
	return e.mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE)
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
		err := e.mw.AutoOrientImage()
		if err != nil {
			return nil, err
		}
	}
	err := e.mw.StripImage()
	if err != nil {
		return nil, err
	}

	return e.mw.GetImageBlob(), nil
}
