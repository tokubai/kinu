package storage
import (
	"io"
	"errors"
	"os"
)

type Storage interface {
	Open() error

	Fetch(key string) ([]byte, error)

	PutFromBlob(key string, image []byte) error
	Put(key string, imageFile io.ReadSeeker) error

	List(key string) ([]StorageItem, error)

	Move(from string, to string) error
}

type StorageItem interface {
	Key() string
	Extension() string
	ImageSize() string
}

var (
	AvailableStorageTypes = []string{ "S3" }
	ErrUnknownStorage = errors.New("specify unknown storage.")
	selectedStorageType string
)

func init() {
	selectedStorageType = os.Getenv("KINU_STORAGE_TYPE")
	if len(selectedStorageType) == 0 {
		panic("must specify KINU_STORAGE_TYPE system environment.")
	}

	var isAvailableStorageType bool
	for _, storageType := range AvailableStorageTypes {
		if selectedStorageType == storageType {
			isAvailableStorageType = true
		}
	}

	if !isAvailableStorageType {
		panic("unknown KINU_STORAGE_TYPE " + selectedStorageType + ".")
	}
}

func Open() (Storage, error) {
	if selectedStorageType == "S3" {
		return openS3Storage()
	} else {
		return nil, ErrUnknownStorage
	}
}
