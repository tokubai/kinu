package resizer
import (
	"github.com/TakatoshiMaeda/kinu/engine"
	"github.com/TakatoshiMaeda/kinu/logger"
)

type ResizeOption struct {
	Width     int
	Height    int
	NeedsAutoCrop bool
	Quality   int
}

const (
	DEFAULT_QUALITY = 70
)

func Run(image []byte, option *ResizeOption) (resizedImage []byte, err error) {
	calculator, err := NewCoodinatesCalculator(option)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	if option.Quality == 0 {
		option.Quality = DEFAULT_QUALITY
	}

	engine, err := engine.New(image)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	engine.SetResizeSize(option.Width, option.Height)

	err = engine.Open()
	if err != nil {
		return nil, logger.ErrorDebug(err)
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
		return nil, logger.ErrorDebug(err)
	}

	if coodinates.CanCrop() {
		err = engine.Crop(coodinates.CropWidth, coodinates.CropHeight, coodinates.WidthOffset, coodinates.HeightOffset)
		if err != nil {
			return nil, logger.ErrorDebug(err)
		}
	}

	return engine.Generate()
}