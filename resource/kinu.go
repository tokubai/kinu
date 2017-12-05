package resource

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/tokubai/kinu/logger"
	"github.com/tokubai/kinu/resizer"
	"github.com/tokubai/kinu/storage"
	"github.com/tokubai/kinu/uploader"
)

var (
	kinuImageFilePathRegexp *regexp.Regexp
)

func init() {
	kinuImageFilePathRegexp = regexp.MustCompile(`(.*).(original|[0-9]{4,5}).kinu\z`)
}

type KinuResource struct {
	Resource

	Category string
	Id       string
}

func (r *KinuResource) FilePath(size string) string {
	return fmt.Sprintf("%s/%s.%s.kinu", r.BasePath(), r.Id, size)
}

func (r *KinuResource) BasePath() string {
	return fmt.Sprintf("%s/%s", r.Category, r.Id)
}

func (r *KinuResource) Fetch(geo *resizer.Geometry) (*Image, error) {
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
		logger.ErrorDebug(err)
	}

	image.Width, err = strconv.Atoi(obj.Metadata["Width"])
	if err != nil {
		logger.ErrorDebug(err)
	}

	return image, nil
}

func (r *KinuResource) MoveTo(category, id string) error {
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

			if kinuImageFilePathRegexp.MatchString(item.Key()) {
				err = st.Move(item.Key(), moveToResource.FilePath(item.ImageSize()))
			} else {
				err = st.Move(item.Key(), moveToResource.BasePath()+"/"+item.Filename())
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

func (r *KinuResource) Store(file io.ReadSeeker) error {
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return &ErrStore{Message: "invalid file"}
	}

	ext := ""
	contentType := http.DetectContentType(imageData)
	switch http.DetectContentType(imageData) {
	case "image/jpeg":
		ext = "jpg"
	case "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	default:
		return &ErrStore{Message: "unsupported filetype, supported jpg or png or gif"}
	}

	uploaders := make([]uploader.Uploader, 0)
	for _, size := range resizer.MiddleImageSizes {
		uploader := &uploader.ImageUploader{
			ImageBlob:   imageData,
			Path:        r.FilePath(size),
			UploadSize:  size,
			ContentType: contentType,
			Ext:         ext,
		}
		uploaders = append(uploaders, uploader)
	}
	uploaders = append(uploaders,
		&uploader.TextFileUploader{
			Path: fmt.Sprintf("%s/filetype.%s", r.BasePath(), ext),
		},
	)

	return uploader.Upload(uploaders)
}
