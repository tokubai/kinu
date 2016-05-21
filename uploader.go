package main

import (
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/TakatoshiMaeda/kinu/storage"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type ErrImageUpload struct {
	error
	Errors []error
}

func (e *ErrImageUpload) Error() string {
	messages := "Image upload error. cause, "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i+1) + ". " + err.Error() + "  "
	}
	return messages
}

var (
	imageUploadSizes = []string{"original", "1000", "2000", "3000"}
)

func UploadImage(imageType string, imageId string, imageFile io.ReadSeeker) error {
	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return &ErrInvalidRequest{Message: "invalid file"}
	}

	contentType := ""
	switch http.DetectContentType(imageData) {
	case "image/jpeg":
		contentType = "jpg"
	case "image/jpg":
		contentType = "jpg"
	default:
		return &ErrInvalidRequest{Message: "unsupported filetype, only support jpg"}
	}

	imageMetadata := NewImageMetadata(imageType, imageId)

	uploaders := make([]Uploader, 0)
	for _, size := range imageUploadSizes {
		uploader := &ImageUploader{
			ImageMetadata: imageMetadata,
			ImageBlob:     imageData,
			UploadSize:    size,
		}
		uploaders = append(uploaders, uploader)
	}

	uploader := &FileTypeTextUploader{
		Filetype:      contentType,
		ImageMetadata: imageMetadata,
	}
	uploaders = append(uploaders, uploader)

	wg := sync.WaitGroup{}
	errs := make(chan error, len(uploaders))
	for _, uploader := range uploaders {
		wg.Add(1)
		go func(u Uploader, errs chan error) {
			defer wg.Done()
			errs <- u.Exec()
		}(uploader, errs)
	}

	wg.Wait()

	close(errs)

	errors := make([]error, 0)
	for err := range errs {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	} else {
		return &ErrImageUpload{Errors: errors}
	}

	return nil
}

type Uploader interface {
	Exec() error
}

type ImageUploader struct {
	Uploader

	ImageMetadata *ImageMetadata
	ImageBlob     []byte
	UploadSize    string
}

func (u *ImageUploader) NeedsResize() bool {
	return u.UploadSize != "original"
}

func (u *ImageUploader) BuildResizeOption() (*resizer.ResizeOption, error) {
	if u.UploadSize == "original" {
		return &resizer.ResizeOption{}, nil
	}

	size, err := strconv.Atoi(u.UploadSize)
	if err != nil {
		return nil, err
	}

	return &resizer.ResizeOption{Width: size, Height: size}, nil
}

func (u *ImageUploader) Exec() error {
	if u.NeedsResize() {
		resizeOption, err := u.BuildResizeOption()
		if err != nil {
			return err
		}

		u.ImageBlob, err = resizer.Run(u.ImageBlob, resizeOption)
		if err != nil {
			return err
		}
	}

	storage, err := storage.Open()
	if err != nil {
		return err
	}

	return storage.PutFromBlob(u.ImageMetadata.FilePath(u.UploadSize), u.ImageBlob)
}

type FileTypeTextUploader struct {
	Uploader

	Filetype      string
	ImageMetadata *ImageMetadata
}

func (u *FileTypeTextUploader) Exec() error {
	storage, err := storage.Open()
	if err != nil {
		return logger.ErrorDebug(err)
	}

	return storage.PutFromBlob(u.ImageMetadata.BasePath()+"/filetype."+u.Filetype, []byte{})
}
