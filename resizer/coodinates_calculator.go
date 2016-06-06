package resizer

import(
	"math"
	"fmt"
)

type CoodinatesCalculator struct {
	ImageWidth  int
	ImageHeight int
	Width       int
	Height      int
}

type ErrInvalidOption struct {
	error
	Message string
}

func (e *ErrInvalidOption) Error() string { return e.Message }

func NewCoodinatesCalculator(option *ResizeOption) (*CoodinatesCalculator, error) {
	if option.Width <= 0 && option.Height <= 0 {
		return nil, &ErrInvalidOption{Message: "option must specify Width or Height"}
	}
	return &CoodinatesCalculator{Width: option.Width, Height: option.Height}, nil
}

func (c *CoodinatesCalculator) SetImageSize(width int, height int) {
	c.ImageWidth = width
	c.ImageHeight = height
}

func (c *CoodinatesCalculator) Resize() (coodinates *Coodinates) {
	coodinates = &Coodinates{}
	switch {
	case c.Width > 0 && c.Height == 0: // Fixed Width
		coodinates.ResizeWidth = c.Width
		coodinates.ResizeHeight = int(float64(c.ImageHeight) * (float64(c.Width) / float64(c.ImageWidth)))
	case c.Width == 0 && c.Height > 0: // Fixed Height
		coodinates.ResizeWidth = int(float64(c.ImageWidth) * (float64(c.Height) / float64(c.ImageHeight)))
		coodinates.ResizeHeight = c.Height
	default: // Fixed Width and Height
		scaleRatio := math.Min(float64(c.Height)/float64(c.ImageHeight), float64(c.Width)/float64(c.ImageWidth))
		coodinates.ResizeWidth = int(float64(c.ImageWidth) * scaleRatio)
		coodinates.ResizeHeight = int(float64(c.ImageHeight) * scaleRatio)
	}
	return coodinates
}

func (c *CoodinatesCalculator) AutoCrop() (coodinates *Coodinates) {
	coodinates = &Coodinates{CropHeight: c.Height, CropWidth: c.Width}

	heightScaleRatio := float64(c.Height) / float64(c.ImageHeight)
	widthScaleRatio := float64(c.Width) / float64(c.ImageWidth)

	scaleRatio := math.Max(heightScaleRatio, widthScaleRatio)

	if heightScaleRatio > widthScaleRatio {
		coodinates.WidthOffset = int((float64(c.ImageWidth)*scaleRatio - float64(c.Width)) / float64(2.0))
		coodinates.ResizeHeight = c.Height
		coodinates.ResizeWidth = int(float64(c.ImageWidth) * scaleRatio)
	} else {
		coodinates.HeightOffset = int((float64(c.ImageHeight)*scaleRatio - float64(c.Height)) / float64(2.0))
		coodinates.ResizeHeight = int(float64(c.ImageHeight) * scaleRatio)
		coodinates.ResizeWidth = c.Width
	}
	return coodinates
}

func (c *CoodinatesCalculator) ManualCrop(option *ResizeOption) (coodinates *Coodinates) {
	coodinates = &Coodinates{}

	assumeRatio := float64(c.ImageWidth) / float64(option.AssumptionWidth)

	coodinates.CropHeight = int(float64(option.CropHeight) * assumeRatio)
	coodinates.CropWidth = int(float64(option.CropWidth) * assumeRatio)
	coodinates.WidthOffset = int(float64(option.CropWidthOffset) * assumeRatio)
	coodinates.HeightOffset = int(float64(option.CropHeightOffset) * assumeRatio)

	coodinates.ResizeWidth = option.Width
	coodinates.ResizeHeight = option.Height
	return coodinates
}

func (c *CoodinatesCalculator) Calc(option *ResizeOption) (coodinates *Coodinates) {
	if option.NeedsAutoCrop {
		return c.AutoCrop()
	} else if option.NeedsManualCrop {
		return c.ManualCrop(option)
	} else {
		return c.Resize()
	}
}

type Coodinates struct {
	ResizeWidth, ResizeHeight int
	CropWidth,   CropHeight   int
	WidthOffset, HeightOffset int
}

func (c *Coodinates) Valid() bool {
	return c.ResizeHeight > 0 && c.ResizeWidth > 0
}

func (c *Coodinates) CanCrop() bool {
	return c.CropWidth > 0 && c.CropHeight > 0
}

func (c *Coodinates) ToString() string {
	return fmt.Sprintf("ResizeWidth: %d, ResizeHeight: %d, CropWidth: %d, CropHeight: %d, WidthOffset: %d HeightOffset: %d",
		c.ResizeWidth,
		c.ResizeHeight,
		c.CropWidth,
		c.CropHeight,
		c.WidthOffset,
		c.HeightOffset,
	)
}
