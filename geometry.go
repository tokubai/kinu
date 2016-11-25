package main

import (
	"fmt"
	"github.com/tokubai/kinu/resizer"
	"strconv"
	"strings"
)

const DEFAULT_QUALITY = 80
const MAX_QUALITY = 100
const MIN_QUALITY = 0

type Geometry struct {
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	Quality            int    `json:"quality"`
	NeedsAutoCrop      bool   `json:"needs_auto_crop"`
	NeedsManualCrop    bool   `json:"needs_manual_crop"`
	CropWidthOffset    int    `json:"cropWidthOffset"`
	CropHeightOffset   int    `json:"cropHeightOffset"`
	CropWidth          int    `json:"cropWidth"`
	CropHeight         int    `json:"cropHeight"`
	AssumptionWidth    int    `json:"assumptionWidth"`
	NeedsOriginalImage bool   `json:"needs_original_image"`
	MiddleImageSize    string `json:"middle_image_size"`
}

const (
	AUTO_CROP = iota
	NORMAL_RESIZE
	ORIGINAL
)

const (
	GEO_NONE = iota
	GEO_WIDTH
	GEO_HEIGHT
	GEO_QUALITY
	GEO_AUTO_CROP
	GEO_MANUAL_CROP
	GEO_WIDTH_OFFSET
	GEO_HEIGHT_OFFSET
	GEO_CROP_WIDTH
	GEO_CROP_HEIGHT
	GEO_ASSUMPTION_WIDTH
	GEO_ORIGINAL
	GEO_MIDDLE
)

func ParseGeometry(geo string) (*Geometry, error) {
	conditions := strings.Split(geo, ",")

	var width, height, quality int
	var middleImageSize = ""
	var pos = GEO_NONE
	var needsAutoCrop, needsManualCrop, needsOriginal bool
	var cropWidthOffset, cropHeightOffset, cropWidth, cropHeight, assumptionWidth int
	for _, condition := range conditions {
		cond := strings.Split(condition, "=")

		if len(cond) < 2 {
			return nil, &ErrInvalidRequest{Message: "invalid geometry, support geometry pattern is key=value,key2=value."}
		}

		switch cond[0] {
		case "w":
			if pos >= GEO_WIDTH {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry w must be fixed order."}
			}
			pos = GEO_WIDTH
			if w, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry w is must be numeric."}
			} else {
				width = w
			}
		case "h":
			if pos >= GEO_HEIGHT {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry h must be fixed order."}
			}
			pos = GEO_HEIGHT
			if h, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry h is must be numeric."}
			} else {
				height = h
			}
		case "q":
			if pos >= GEO_QUALITY {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry q must be fixed order."}
			}
			pos = GEO_QUALITY
			if q, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry q is must be numeric."}
			} else if q > MAX_QUALITY || q < MIN_QUALITY {
				return nil, &ErrInvalidRequest{Message: "q is under " + strconv.Itoa(MAX_QUALITY) + " and over " + strconv.Itoa(MIN_QUALITY)}
			} else {
				quality = q
			}
		case "c":
			if pos >= GEO_AUTO_CROP {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry c must be fixed order."}
			}
			pos = GEO_AUTO_CROP
			if cond[1] == "true" {
				needsAutoCrop = true
			} else {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry c must be true or manual."}
			}
		case "mc":
			if pos >= GEO_MANUAL_CROP {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry mc must be fixed order."}
			}
			pos = GEO_MANUAL_CROP
			if cond[1] == "true" {
				needsManualCrop = true
			} else {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry mc must be true or manual."}
			}
		case "wo":
			if pos >= GEO_WIDTH_OFFSET {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry wo must be fixed order."}
			}
			pos = GEO_WIDTH_OFFSET
			if wo, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry wo is must be numeric."}
			} else {
				cropWidthOffset = wo
			}
		case "ho":
			if pos >= GEO_HEIGHT_OFFSET {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry ho must be fixed order."}
			}
			pos = GEO_HEIGHT_OFFSET
			if ho, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry ho is must be numeric."}
			} else {
				cropHeightOffset = ho
			}
		case "cw":
			if pos >= GEO_CROP_WIDTH {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry cw must be fixed order."}
			}
			pos = GEO_CROP_WIDTH
			if cw, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry cw is must be numeric."}
			} else {
				cropWidth = cw
			}
		case "ch":
			if pos >= GEO_CROP_HEIGHT {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry ch must be fixed order."}
			}
			pos = GEO_CROP_HEIGHT
			if ch, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry ch is must be numeric."}
			} else {
				cropHeight = ch
			}
		case "aw":
			if pos >= GEO_ASSUMPTION_WIDTH {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry as must be fixed order."}
			}
			pos = GEO_ASSUMPTION_WIDTH
			if aw, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry as is must be numeric."}
			} else {
				assumptionWidth = aw
			}
		case "o":
			if pos >= GEO_ORIGINAL {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry o must be fixed order."}
			}
			pos = GEO_ORIGINAL
			if cond[1] == "true" {
				needsOriginal = true
			} else {
				needsOriginal = false
			}
		case "m":
			if pos >= GEO_MIDDLE {
				return nil, &ErrInvalidGeometryOrderRequest{Message: "geometry m must be fixed order."}
			}
			pos = GEO_MIDDLE
			if cond[1] == "true" {
				middleImageSize = "1000"
			} else {
				for _, size := range middleImageSizes {
					if cond[1] == size {
						middleImageSize = cond[1]
						break
					}
				}
			}
			if len(middleImageSize) == 0 {
				return nil, &ErrInvalidRequest{Message: "must specify valid middle image size."}
			}
		}
	}

	if len(middleImageSize) == 0 && width == 0 && height == 0 && needsOriginal == false {
		return nil, &ErrInvalidRequest{Message: "must specify width or height when not original mode."}
	}

	if needsManualCrop && (cropWidth == 0 || cropHeight == 0 || assumptionWidth == 0) {
		return nil, &ErrInvalidRequest{Message: "must specify crop width, crop height and assumption width when manual crop mode."}
	}

	if quality == 0 {
		quality = DEFAULT_QUALITY
	}

	return &Geometry{
		Width: width, Height: height,
		Quality: quality,
		NeedsAutoCrop: needsAutoCrop,
		NeedsManualCrop: needsManualCrop,
		CropWidthOffset: cropWidthOffset,
		CropHeightOffset: cropHeightOffset,
		CropWidth: cropWidth,
		CropHeight: cropHeight,
		AssumptionWidth: assumptionWidth,
		MiddleImageSize: middleImageSize,
		NeedsOriginalImage: needsOriginal}, nil
}

func (g *Geometry) ResizeMode() int {
	if g.NeedsAutoCrop {
		return AUTO_CROP
	}

	if g.NeedsOriginalImage {
		return ORIGINAL
	}

	return NORMAL_RESIZE
}

func (g *Geometry) ToResizeOption() (resizeOption *resizer.ResizeOption) {
	return &resizer.ResizeOption{
		Width:            g.Width,
		Height:           g.Height,
		Quality:          g.Quality,
		NeedsAutoCrop:    g.NeedsAutoCrop,
		NeedsManualCrop:  g.NeedsManualCrop,
		CropWidthOffset:  g.CropWidthOffset,
		CropHeightOffset: g.CropHeightOffset,
		CropWidth:        g.CropWidth,
		CropHeight:       g.CropHeight,
		AssumptionWidth:  g.AssumptionWidth,
	}
}

func (g *Geometry) ToString() string {
	return fmt.Sprintf("Width: %d, Height: %d, Quality: %d, NeedsAutoCrop: %t, NeedsManualCrop: %t, NeedsOriginalImage: %t", g.Width, g.Height, g.Quality, g.NeedsAutoCrop, g.NeedsManualCrop, g.NeedsOriginalImage)
}
