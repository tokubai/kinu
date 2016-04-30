package main

import (
	"strconv"
	"strings"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"fmt"
)

const DEFAULT_QUALITY = 80
const MAX_QUALITY = 100
const MIN_QUALITY = 0

type Geometry struct {
	Width      int `json:"width"`
	Height     int `json:"height"`
	Quality    int `json:"quality"`
	NeedsAutoCrop bool `json:"needs_auto_crop"`
	NeedsOriginalImage bool `json:"needs_original_image"`
}

const (
	AUTO_CROP = iota
	NORMAL_RESIZE
	ORIGINAL
)

func ParseGeometry(geo string) (*Geometry, error) {
	conditions := strings.Split(geo, ",")

	var width, height, quality int
	var needsAutoCrop, needsOriginal bool
	for _, condition := range conditions {
		cond := strings.Split(condition, "=")

		if len(cond) < 2 {
			return nil, &ErrInvalidRequest{Message: "invalid geometry, support geometry pattern is key=value,key2=value."}
		}

		switch cond[0] {
		case "w":
			if w, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry w is must be numeric."}
			} else {
				width = w
			}
		case "h":
			if h, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry h is must be numeric."}
			} else {
				height = h
			}
		case "q":
			if q, err := strconv.Atoi(cond[1]); err != nil {
				return nil, &ErrInvalidRequest{Message: "geometry q is must be numeric."}
			} else if q > MAX_QUALITY || q < MIN_QUALITY {
				return nil, &ErrInvalidRequest{Message: "q is under " + strconv.Itoa(MAX_QUALITY) + " and over " + strconv.Itoa(MIN_QUALITY)}
			} else {
				quality = q
			}
		case "c":
			if cond[1] == "true" {
				needsAutoCrop = true
			} else {
				needsAutoCrop = false
			}
		case "o":
			if cond[1] == "true" {
				needsOriginal = true
			} else {
				needsOriginal = false
			}
		}
	}

	if width == 0 && height == 0 && needsOriginal == false {
		return nil, &ErrInvalidRequest{Message: "must specify width or height when not original mode."}
	}

	if quality == 0 {
		quality = DEFAULT_QUALITY
	}

	return &Geometry{Width: width, Height: height, Quality: quality, NeedsAutoCrop: needsAutoCrop, NeedsOriginalImage: needsOriginal}, nil
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
		Width: g.Width,
		Height: g.Height,
		Quality: g.Quality,
		NeedsAutoCrop: g.NeedsAutoCrop,
	}
}

func (g *Geometry) ToString() string {
	return fmt.Sprintf("Width: %d, Height: %d, Quality: %d, NeedsAutoCrop: %t, NeedsOriginalImage: %t", g.Width, g.Height, g.Quality, g.NeedsAutoCrop, g.NeedsOriginalImage)
}