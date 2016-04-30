package main
import (
	"io"
	"io/ioutil"
	"net/http"
	"github.com/TakatoshiMaeda/kinu/resizer"
	"strconv"
	"github.com/TakatoshiMaeda/kinu/storage"
	"sync"
)

type ErrImageUpload struct {
	error
	Errors []error
}

func (e *ErrImageUpload) Error() string {
	messages := "Image upload error.\ncause "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i) + ". " + err.Error() + "\n"
	}
	return messages
}

var (
	imageUploadSizes = []string{ "original", "1000", "2000", "3000" }
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

	uploaders := make([]*Uploader, 0)
	for _, size := range imageUploadSizes {
		uploader := &Uploader{
			ImageMetadata: NewImageMetadata(imageType, imageId),
			ImageBlob: imageData,
			Extension: contentType,
			UploadSize: size,
		}
		uploaders = append(uploaders, uploader)
	}

	// FIXME: ErrorHandling
	wg := sync.WaitGroup{}
	// errs := make(chan error)
	for _, uploader := range uploaders {
		wg.Add(1)
		//go func(u *Uploader, errs chan error) {
		go func(u *Uploader) {
			defer wg.Done()
			// errs <- u.Exec()
			u.Exec()
		//}(uploader, errs)
		}(uploader)
	}
	wg.Wait()

	// close(errs)

	// errors := make([]error, 0)
	// for err := range errs {
	// 	if err != nil {
	// 		errors = append(errors, err)
	// 	}
	// }

	// if len(errors) == 0 {
	// 	return nil
	// } else {
	// 	return &ErrImageUpload{Errors: errors}
	// }
	return nil
}

type Uploader struct {
	ImageMetadata *ImageMetadata
	ImageBlob []byte
	Extension string
	UploadSize string
}

func (u *Uploader) NeedsResize() bool {
	return u.UploadSize != "original"
}

func (u *Uploader) BuildResizeOption() (*resizer.ResizeOption, error) {
	if u.UploadSize == "original" {
		return &resizer.ResizeOption{}, nil
	}

	size, err := strconv.Atoi(u.UploadSize)
	if err != nil {
		return nil, err
	}

	return &resizer.ResizeOption{Width: size, Height: size}, nil
}

func (u *Uploader) Exec() error {
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

	return storage.PutFromBlob(u.ImageMetadata.FilePath(u.UploadSize, u.Extension), u.ImageBlob)
}