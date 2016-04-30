package main
import (
	"strings"
	"path/filepath"
	"fmt"
)

type ImageMetadata struct {
	ImageType string
	Id string
}

var (
	validExtensions = []string{ "jpg", "jpeg" }
)

func ExtractExtension(filename string) string {
	return strings.Replace(filepath.Ext(filename), ".", "", 1)
}

func IsValidImageExt(ext string) bool {
	for _, e := range validExtensions {
		if e == ext {
			return true
		}
	}
	return false
}

func NewImageMetadata(imageType string, id string) (*ImageMetadata) {
	return &ImageMetadata{
		ImageType: imageType,
		Id: id,
	}
}

func (i *ImageMetadata) FileMiddleImagePath(ext string) string {
	return i.FilePath("1000", ext)
}

func (i *ImageMetadata) FileOriginPath(ext string) string {
	return i.FilePath("origin", ext)
}

func (i *ImageMetadata) FilePath(size string, ext string) string {
	return fmt.Sprintf("%s/%s.%s.%s", i.BasePath(), i.Id, size, ext)
}

func (i *ImageMetadata) BasePath() string {
	return fmt.Sprintf("%s/%s", i.ImageType, i.Id)
}
