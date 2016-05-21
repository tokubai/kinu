package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/TakatoshiMaeda/kinu/storage"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/Sirupsen/logrus"
)

type ImageGetRequest struct {
	imageMetadata *ImageMetadata

	extension string
	geometry *Geometry
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

	ext := ExtractExtension(filename)
	if len(ext) == 0 {
		return nil, &ErrInvalidRequest{Message: "invalid file extension"}
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
		"geometry": geometry.ToString(),
		"image_type": imageType,
		"image_id": id,
		"extension": ext,
	}).Debug("parse success image get request.")

	return &ImageGetRequest{imageMetadata: NewImageMetadata(imageType, id), geometry: geometry, extension: ext}, nil
}

func (r *ImageGetRequest) FetchImage() (image []byte, err error) {
	storage, err := storage.Open()
	if err != nil {
		return nil, err
	}
	return storage.Fetch(r.imageMetadata.FileMiddleImagePath())
}

func (r *ImageGetRequest) NeedsOriginalImage() bool {
	return r.geometry.NeedsOriginalImage
}

func (r *ImageGetRequest) ToResizeOption() (*resizer.ResizeOption) {
	return r.geometry.ToResizeOption()
}
