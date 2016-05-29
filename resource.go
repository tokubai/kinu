package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"fmt"
	"github.com/TakatoshiMaeda/kinu/storage"
	"github.com/TakatoshiMaeda/kinu/logger"
	"strconv"
	"regexp"
	"github.com/Sirupsen/logrus"
)

var (
	validExtensions = []string{"jpg", "jpeg"}
	middleImageSizes = []string{"original", "1000", "2000", "3000"}
	imageFilePathRegexp *regexp.Regexp
)

func init() {
	imageFilePathRegexp = regexp.MustCompile(`(.*).(original|[0-9]{4,5}).kinu\z`)
}

type Resource struct {
	Category string
	Id       string
}

type Image struct {
	Width int
	Height int
	Body []byte
}

type ErrMove struct {
	error
	Errors []error
}

func (e *ErrMove) Error() string {
	messages := "Move error. cause, "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i+1) + ". " + err.Error() + "  "
	}
	return messages
}

func NewResource(category string, id string) *Resource {
	return &Resource{
		Category: category,
		Id: id,
	}
}

func (r *Resource) FilePath(size string) string {
	return fmt.Sprintf("%s/%s.%s.kinu", r.BasePath(), r.Id, size)
}

func (r *Resource) BasePath() string {
	return fmt.Sprintf("%s/%s", r.Category, r.Id)
}

func (r *Resource) Fetch(geo *Geometry) (*Image, error) {
	var middleImageSize string
	if geo.NeedsOriginalImage {
		middleImageSize = "original"
	} else if geo.Height <= 1000 && geo.Width <= 1000 {
		middleImageSize = "1000"
	} else if geo.Height <= 2000 && geo.Width <= 2000 {
		middleImageSize = "2000"
	} else if geo.Height <= 3000 && geo.Width <= 3000 {
		middleImageSize = "3000"
	} else {
		middleImageSize = "original"
	}

	image := &Image{}

	st, err := storage.Open()
	if err != nil {
		return image, logger.ErrorDebug(err)
	}

	obj, err := st.Fetch(r.FilePath(middleImageSize))
	if err != nil {
		return image, logger.ErrorDebug(err)
	}

	image.Body = obj.Body

	logger.WithFields(logrus.Fields{
		"metadata": obj.Metadata,
	}).Debug("metadata")

	image.Height, err = strconv.Atoi(obj.Metadata["Height"])
	if err != nil {
		return image, logger.ErrorDebug(err)
	}

	image.Width, err = strconv.Atoi(obj.Metadata["Width"])
	if err != nil {
		return image, logger.ErrorDebug(err)
	}

	return image, nil
}

func (r *Resource) MoveTo(category, id string) error {
	st, err := storage.Open()
	if err != nil {
		return logger.ErrorDebug(err)
	}

	items, err := st.List(r.BasePath())
	if err != nil {
		return logger.ErrorDebug(err)
	}

	moveToResource := NewResource(category, id)

	wg := sync.WaitGroup{}
	errs := make(chan error, len(items))
	for _, item := range items {
		wg.Add(1)
		go func(item storage.StorageItem) {
			defer wg.Done()
			st, err := storage.Open()
			if err != nil {
				errs <- logger.ErrorDebug(err)
				return
			}

			if imageFilePathRegexp.MatchString(item.Key()) {
				err = st.Move(item.Key(), moveToResource.FilePath(item.ImageSize()))
			} else {
				err = st.Move(item.Key(), moveToResource.BasePath() + "/" + item.Filename())
			}

			if err != nil {
				errs <- logger.ErrorDebug(err)
				return
			}

			errs <- nil
		}(item)
	}
	wg.Wait()

	close(errs)

	errors := make([]error, 0)
	for err := range errs {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		return &ErrAttachFromSandbox{Errors: errors}
	}

	return nil
}

func (r *Resource) Store(file io.ReadSeeker) error {
	imageData, err := ioutil.ReadAll(file)
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

	uploaders := make([]Uploader, 0)
	for _, size := range middleImageSizes {
		uploader := &ImageUploader{
			ImageBlob:     imageData,
			Path: r.FilePath(size),
			UploadSize:    size,
		}
		uploaders = append(uploaders, uploader)
	}
	uploaders = append(uploaders,
		&TextFileUploader{
			Path: fmt.Sprintf("%s/filetype.%s", r.BasePath(), contentType),
		},
	)

	return Upload(uploaders)
}
