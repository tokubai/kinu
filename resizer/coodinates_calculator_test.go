package resizer

import (
	"testing"
)

func TestNewCoodinatesCalculator(t *testing.T) {
	_, err := NewCoodinatesCalculator(&ResizeOption{})
	if err == nil {
		t.Errorf("must returns error when height and width option specify missing ")
	}

	_, err = NewCoodinatesCalculator(&ResizeOption{Width: 100})
	if err != nil {
		t.Errorf("height or width option specify when should be success")
	}

	_, err = NewCoodinatesCalculator(&ResizeOption{Height: 100})
	if err != nil {
		t.Errorf("height or width option specify when should be success")
	}

	_, err = NewCoodinatesCalculator(&ResizeOption{Height: 100, Width: 100})
	if err != nil {
		t.Errorf("height or width option specify when should be success")
	}
}

func TestCoodinatesCalculator_Resize(t *testing.T) {
	calculator := CoodinatesCalculator{Width: 100, Height: 50}

	calculator.SetImageSize(200, 300)
	coodinates := calculator.Resize()
	if coodinates.ResizeWidth == 33 && coodinates.ResizeHeight == 50 {
		// Success
	} else {
		t.Errorf(`resize coodinates must fit of bigger scale change
		expect -> Width: 33, Height: 50 / actual -> Width: %d, Height: %d`,
			coodinates.ResizeWidth, coodinates.ResizeHeight,
		)
	}

	calculator.SetImageSize(500, 200)
	coodinates = calculator.Resize()
	if coodinates.ResizeWidth == 100 && coodinates.ResizeHeight == 40 {
		// Success
	} else {
		t.Errorf(`resize coodinates must fit of bigger scale change
		expect -> Width: 100, Height: 40 / actual -> Width: %d, Height: %d`,
			coodinates.ResizeWidth, coodinates.ResizeHeight,
		)
	}

	calculator = CoodinatesCalculator{Width: 100}
	calculator.SetImageSize(500, 200)
	coodinates = calculator.Resize()
	if coodinates.ResizeWidth == 100 && coodinates.ResizeHeight == 40 {
		// Success
	} else {
		t.Errorf(`resize coodinates must fit of specify distance,
		expect -> Width: 100, Height: 40 / actual -> Width: %d, Height: %d`,
			coodinates.ResizeWidth, coodinates.ResizeHeight,
		)
	}

	calculator = CoodinatesCalculator{Height: 50}
	calculator.SetImageSize(500, 200)
	coodinates = calculator.Resize()
	if coodinates.ResizeWidth == 125 && coodinates.ResizeHeight == 50 {
		// Success
	} else {
		t.Errorf(`resize coodinates must fit of specify distance,
		expect -> Width: 125, Height: 50 / actual -> Width: %d, Height: %d`,
			coodinates.ResizeWidth, coodinates.ResizeHeight,
		)
	}
}

func TestCoodinatesCalculator_AutoCrop(t *testing.T) {
	calculator := CoodinatesCalculator{Width: 100, Height: 50}

	calculator.SetImageSize(200, 300)
	c := calculator.AutoCrop()
	if c.CropWidth == 100 && c.CropHeight == 50 &&
		c.ResizeWidth == 100 && c.ResizeHeight == 150 &&
		c.WidthOffset == 0 && c.HeightOffset == 50 {
		// Success
	} else {
		t.Errorf(`auto crop coodinates must fit of specify image size and has offset smaller scale change,
		expect -> CropWidth: 100, CropHeight: 50, Width: 125, Height: 50, WidthOffset: 0, HeightOffset: 50
		actual -> CropWidth: %d, CropHeight: %d, Width: %d, Height: %d, WidthOffset: %d, HeightOffset: %d`,
			c.CropWidth, c.CropHeight, c.ResizeWidth, c.ResizeHeight, c.WidthOffset, c.HeightOffset,
		)
	}

	calculator.SetImageSize(500, 200)
	c = calculator.AutoCrop()
	if c.CropWidth == 100 && c.CropHeight == 50 &&
		c.ResizeWidth == 125 && c.ResizeHeight == 50 &&
		c.WidthOffset == 12 && c.HeightOffset == 0 {
		// Success
	} else {
		t.Errorf(`auto crop coodinates must fit of specify image size and has offset smaller scale change,
		expect -> CropWidth: 100, CropHeight: 50, Width: 125, Height: 50, WidthOffset: 12, HeightOffset: 0
		actual -> CropWidth: %d, CropHeight: %d, Width: %d, Height: %d, WidthOffset: %d, HeightOffset: %d`,
			c.CropWidth, c.CropHeight, c.ResizeWidth, c.ResizeHeight, c.WidthOffset, c.HeightOffset,
		)
	}
}

func TestCoodinates_CanCrop(t *testing.T) {
	c := &Coodinates{CropWidth: 100, CropHeight: 50}
	if !c.CanCrop() {
		t.Errorf("CropWidth and CropHeight specify when actual can crop")
	}

	c = &Coodinates{CropWidth: 100}
	if c.CanCrop() {
		t.Errorf("CropHeight not set when not can crop")
	}

	c = &Coodinates{CropHeight: 50}
	if c.CanCrop() {
		t.Errorf("CropWidth not set when not can crop")
	}
}