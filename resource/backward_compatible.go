package resource

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tokubai/kinu/logger"
	"github.com/tokubai/kinu/resizer"
	"github.com/tokubai/kinu/storage"
	"github.com/tokubai/kinu/uploader"
)

var (
	imageFilePathRegexp *regexp.Regexp
)

type BackwardCompatibleResource struct {
	Resource

	Category string
	Id       string
}

var (
	ErrOriginalImageNotFound = errors.New("not found original image")
)

func (r *BackwardCompatibleResource) FilePath(size string) string {
	if size == "1000" {
		return fmt.Sprintf("%s/%s.jpg", r.BasePath(), r.Id)
	} else {
		return fmt.Sprintf("%s/%s.original.%s.jpg", r.BasePath(), r.Id, time.Now().Format("20060102150405"))
	}
}

func (r *BackwardCompatibleResource) RecentOriginalFileKey() (string, error) {
	s, err := storage.Open()
	if err != nil {
		return "", logger.ErrorDebug(err)
	}

	items, err := s.List(fmt.Sprintf("%s/%s/%s.original", r.Category, r.Id, r.Id))
	if err != nil {
		return "", logger.ErrorDebug(err)
	}

	if len(items) == 0 {
		return "", ErrOriginalImageNotFound
	}

	var recentTimestamp int
	var recentItem storage.StorageItem
	for _, i := range items {
		splittedKey := strings.Split(i.Key(), ".")
		timestamp, _ := strconv.Atoi(splittedKey[len(splittedKey)-2])
		if timestamp >= recentTimestamp {
			recentItem = i
			recentTimestamp = recentTimestamp
			logger.WithFields(logrus.Fields{"key": i.Key(), "timestamp": timestamp}).Debug("update recent original image")
		}
	}

	return recentItem.Key(), nil
}

func (r *BackwardCompatibleResource) BasePath() string {
	return fmt.Sprintf("%s/%s", r.Category, r.Id)
}

func (r *BackwardCompatibleResource) Fetch(geo *resizer.Geometry) (*Image, error) {
	var middleImageSize string
	if geo.NeedsOriginalImage {
		middleImageSize = "original"
	} else if len(geo.MiddleImageSize) != 0 {
		middleImageSize = geo.MiddleImageSize
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

	var path string
	if middleImageSize == "1000" {
		path = r.FilePath(middleImageSize)
	} else {
		path, err = r.RecentOriginalFileKey()
		if err != nil {
			// There are cases where there is no original image and only an intermediate image exists.
			path = r.FilePath("1000")
		}
	}

	obj, err := st.Fetch(path)
	if err != nil {
		return image, logger.ErrorDebug(err)
	}

	image.Body = obj.Body

	logger.WithFields(logrus.Fields{
		"metadata": obj.Metadata,
	}).Debug("metadata")

	image.Height, err = strconv.Atoi(obj.Metadata["Height"])
	if err != nil {
		logger.ErrorDebug(err)
	}

	image.Width, err = strconv.Atoi(obj.Metadata["Width"])
	if err != nil {
		logger.ErrorDebug(err)
	}

	image.ContentType = obj.Metadata["Content-Type"]

	return image, nil
}

func (r *BackwardCompatibleResource) MoveTo(category, id string) error {
	st, err := storage.Open()
	if err != nil {
		return logger.ErrorDebug(err)
	}

	items, err := st.List(r.BasePath())
	if err != nil {
		return logger.ErrorDebug(err)
	}

	moveToResource := New(category, id)

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

			if strings.Contains(item.Key(), "filetype") {
				err = st.Move(item.Key(), moveToResource.BasePath()+"/"+item.Filename())
			} else {
				err = st.Move(item.Key(), moveToResource.FilePath(item.ImageSize()))
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

func (r *BackwardCompatibleResource) Store(file io.ReadSeeker) error {
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return &ErrStore{Message: "invalid file"}
	}

	ext := ""
	contentType := http.DetectContentType(imageData)
	switch contentType {
	case "image/jpeg":
		ext = "jpg"
	case "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "application/pdf":
		ext = "pdf"
	case "image/bmp":
		ext = "bmp"
	default:
		return &ErrStore{Message: "unsupported filetype, supported jpg or png or gif or pdf"}
	}

	uploaders := make([]uploader.Uploader, 0)
	for _, size := range []string{"original", "1000"} {
		uploaders = append(uploaders, &uploader.ImageUploader{
			ImageBlob:   imageData,
			Path:        r.FilePath(size),
			UploadSize:  size,
			ContentType: contentType,
			Ext:         ext,
		})
	}

	uploaders = append(uploaders,
		&uploader.TextFileUploader{
			Path: fmt.Sprintf("%s/filetype.%s", r.BasePath(), ext),
		},
	)

	return uploader.Upload(uploaders)
}
