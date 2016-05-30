package resizer

import (
	"github.com/Sirupsen/logrus"
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

	var coodinates *Coodinates
	if option.HasSizeHint() {
		calculator.SetImageSize(option.SizeHintWidth, option.SizeHintHeight)
		coodinates = calculator.Calc(option)
		engine.SetSizeHint(coodinates.ResizeWidth, coodinates.ResizeHeight)
		logger.WithFields(logrus.Fields{
			"width_size_hint":  option.SizeHintWidth,
			"height_size_hint": option.SizeHintHeight,
		}).Debug("size hint")
	} else {
		logger.WithFields(logrus.Fields{
			"width_size_hint":  option.SizeHintWidth,
			"height_size_hint": option.SizeHintHeight,
		}).Debug("not set size hint")
	}

	err = engine.Open()
	if err != nil {
		return &ResizeResult{err: logger.ErrorDebug(err)}
	}

	defer engine.Close()

	if !coodinates.Valid() {
		calculator.SetImageSize(engine.GetImageWidth(), engine.GetImageHeight())
		coodinates = calculator.Calc(option)
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
