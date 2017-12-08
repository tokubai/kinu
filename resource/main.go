package resource

import (
	"io"
	"strconv"

	"github.com/tokubai/kinu/config"
	"github.com/tokubai/kinu/resizer"
)

var (
	ValidExtensions      = []string{"jpg", "jpeg", "png", "webp", "gif"}
	selectedResourceType string
)

type Resource interface {
	FilePath(size string) string
	BasePath() string
	Fetch(geo *resizer.Geometry) (*Image, error)
	MoveTo(category, id string) error
	Store(file io.ReadSeeker) error
}

type Image struct {
	Width       int
	Height      int
	ContentType string
	Body        []byte
}

type ErrMove struct {
	error
	Errors []error
}

type ErrAttachFromSandbox struct {
	error
	Errors []error
}

func (e *ErrAttachFromSandbox) Error() string {
	messages := "Image attach from sandbox error. cause, "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i+1) + ". " + err.Error() + "  "
	}
	return messages
}

type ErrStore struct {
	error
	Message string
}

func (e *ErrStore) Error() string { return e.Message }

func (e *ErrMove) Error() string {
	messages := "Move error. cause, "
	for i, err := range e.Errors {
		messages = messages + strconv.Itoa(i+1) + ". " + err.Error() + "  "
	}
	return messages
}

func New(category string, id string) Resource {
	if config.BackwardCompatibleMode {
		return &BackwardCompatibleResource{
			Category: category,
			Id:       id,
		}
	} else {
		return &KinuResource{
			Category: category,
			Id:       id,
		}
	}
}
