package engine

import (
	"github.com/disintegration/gift"
	"image"
	_ "image/jpeg"
	"bytes"
	"github.com/TakatoshiMaeda/kinu/logger"
	"image/jpeg"
)

type GiftEngine struct {
	ResizeEngine

	width, height int

	blobImage []byte

	image image.Image
	imageConfig image.Config

	gift *gift.GIFT
}

func newGiftEngine(image []byte) (e *GiftEngine) {
	return &GiftEngine{
		gift: gift.New(),
		blobImage: image,
	}
}

func (e *GiftEngine) SetResizeSize(width int, height int) {  }

func (e *GiftEngine) Open() error {
	buf := bytes.NewBuffer(e.blobImage)
	img, _, err := image.Decode(buf)
	if err != nil {
		return logger.ErrorDebug(err)
	}
	e.image = img

	buf = bytes.NewBuffer(e.blobImage)
	imgConfig, _, err := image.DecodeConfig(buf)
	if err != nil {
		return logger.ErrorDebug(err)
	}
	e.imageConfig = imgConfig

	return nil
}

func (e *GiftEngine) Close() {  }

func (e *GiftEngine) GetImageHeight() int {
	return e.imageConfig.Height
}

func (e *GiftEngine) GetImageWidth() int {
	return e.imageConfig.Width
}

func (e *GiftEngine) Resize(width int, height int) error {
	e.gift.Add(gift.Resize(width, height, gift.LanczosResampling))
	return nil
}

func (e *GiftEngine) Crop(width int, height int, startX int, startY int) error {
	e.gift.Add(gift.Crop(image.Rect(startX, startY, startX + width, startY + height)))
	return nil
}

func (e *GiftEngine) Generate() ([]byte, error) {
	dst := image.NewRGBA(e.gift.Bounds(e.image.Bounds()))
	e.gift.Draw(dst, e.image)

	imageBuf := new(bytes.Buffer)
	err := jpeg.Encode(imageBuf, dst, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}
	return imageBuf.Bytes(), nil
}
