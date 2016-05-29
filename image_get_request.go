package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/julienschmidt/httprouter"
	"path/filepath"
	"strings"
)

type ImageGetRequest struct {
	Category  string
	Id        string
	Geometry  *Geometry
	extension string
}

func ExtractId(filename string) string {
	return strings.Split(filename, ".")[0]
}

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

func NewImageGetRequest(ps httprouter.Params) (*ImageGetRequest, error) {
	imageType := ps.ByName("type")
	if len(imageType) == 0 {
		return nil, &ErrInvalidRequest{Message: "required image type."}
	}

	filename := ps.ByName("filename")
	if len(filename) == 0 {
		return nil, &ErrInvalidRequest{Message: "required filename."}
	}

	id := ExtractId(filename)
	if len(id) == 0 {
		return nil, &ErrInvalidRequest{Message: "invalid filename"}
	}

	geometry, err := ParseGeometry(ps.ByName("geometry"))
	if err != nil {
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"geometry":   geometry.ToString(),
		"image_type": imageType,
		"image_id":   id,
	}).Debug("parse success image get request.")

	return &ImageGetRequest{Category: imageType, Id: id, Geometry: geometry}, nil
}
