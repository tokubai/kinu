package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ImageMetadata struct {
	ImageType string
	Id        string
}

var (
	validExtensions = []string{"jpg", "jpeg"}
)

func ExtractExtension(filename string) string {
	return strings.Replace(filepath.Ext(filename), ".", "", 1)
}

func ExtractId(filename string) string {
	return strings.Split(filename, ".")[0]
}

func IsValidImageExt(ext string) bool {
	for _, e := range validExtensions {
		if e == ext {
			return true
		}
	}
	return false
}

func NewImageMetadata(imageType string, id string) *ImageMetadata {
	return &ImageMetadata{
		ImageType: imageType,
		Id:        id,
	}
}

func (i *ImageMetadata) FileMiddleImagePath() string {
	return i.FilePath("1000")
}

func (i *ImageMetadata) FileOriginPath() string {
	return i.FilePath("origin")
}

// images/1/1.kinu -> 1000x1000 default middle image
// images/1/1.2000.kinu -> 2000x2000 larger middle image
// images/1/1.3000.kinu -> 3000x3000 more larger middle image
// images/1/1.original.kinu -> original image
func (i *ImageMetadata) FilePath(size string) string {
	if size == "1000" || size == "" {
		return fmt.Sprintf("%s/%s.kinu", i.BasePath(), i.Id)
	} else {
		return fmt.Sprintf("%s/%s.%s.kinu", i.BasePath(), i.Id, size)
	}
}

func (i *ImageMetadata) BasePath() string {
	return fmt.Sprintf("%s/%s", i.ImageType, i.Id)
}
