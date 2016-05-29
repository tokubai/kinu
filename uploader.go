package main

import (
)
import (
	"strconv"
	"sync"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"github.com/TakatoshiMaeda/kinu/storage"
	"github.com/TakatoshiMaeda/kinu/logger"
)

type ErrUpload struct {
	error
	Errors []error
}

func (e *ErrUpload) Error() string {
	messages := "Upload error. cause, "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i+1) + ". " + err.Error() + "  "
	}
	return messages
}

func Upload(uploaders []Uploader) error {
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
		return &ErrUpload{Errors: errors}
	}

	return nil
}

type Uploader interface {
	Exec() error
}

type ImageUploader struct {
	Uploader
	Path string
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

	// TODO: SetMetadata

	return storage.PutFromBlob(u.Path, u.ImageBlob)
}

type TextFileUploader struct {
	Uploader

	Body string
	Path string
}

func (u *TextFileUploader) Exec() error {
	storage, err := storage.Open()
	if err != nil {
		return logger.ErrorDebug(err)
	}
	return storage.PutFromBlob(u.Path, []byte(u.Body))
}
