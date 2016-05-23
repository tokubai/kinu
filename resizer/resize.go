package resizer

import (
	"github.com/TakatoshiMaeda/kinu/engine"
	"github.com/TakatoshiMaeda/kinu/logger"
)

func Resize(image []byte, option *ResizeOption) (result *ResizeResult) {
	calculator, err := NewCoodinatesCalculator(option)
	if err != nil {
		return &ResizeResult{err: logger.ErrorDebug(err)}
	}

	if option.Quality == 0 {
		option.Quality = DEFAULT_QUALITY
	}

	engine, err := engine.New(image)
	if err != nil {
		return &ResizeResult{err: logger.ErrorDebug(err)}
	}

	engine.SetResizeSize(option.Width, option.Height)

	err = engine.Open()
	if err != nil {
		return &ResizeResult{err: logger.ErrorDebug(err)}
	}

	defer engine.Close()

	calculator.SetImageSize(engine.GetImageWidth(), engine.GetImageHeight())

	var coodinates *Coodinates
	if option.NeedsAutoCrop {
		coodinates = calculator.AutoCrop()
	} else {
		coodinates = calculator.Resize()
	}

	err = engine.Resize(coodinates.ResizeWidth, coodinates.ResizeHeight)
	if err != nil {
		return &ResizeResult{err: logger.ErrorDebug(err)}
	}

	if coodinates.CanCrop() {
		err = engine.Crop(coodinates.CropWidth, coodinates.CropHeight, coodinates.WidthOffset, coodinates.HeightOffset)
		if err != nil {
			return &ResizeResult{err: logger.ErrorDebug(err)}
		}
	}

	resultImage, err := engine.Generate()
	return &ResizeResult{image: resultImage, err: err}
}
