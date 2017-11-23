package storage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tokubai/kinu/logger"
)

type FileStorage struct {
	Storage

	baseDirectory string
}

type FileStorageItem struct {
	StorageItem

	Name string
}

func openFileStorage() (Storage, error) {
	s := &FileStorage{}
	err := s.Open()
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}
	return s, nil
}

func (s *FileStorage) Open() error {
	s.baseDirectory = os.Getenv("KINU_FILE_DIRECTORY")
	if len(s.baseDirectory) == 0 {
		return &ErrInvalidStorageOption{Message: "KINU_FILE_DIRECTORY system env is required"}
	}

	logger.WithFields(logrus.Fields{
		"base_directory": s.baseDirectory,
	}).Debug("open file storage")

	return nil
}

func (s *FileStorage) BuildKey(key string) string {
	return s.baseDirectory + "/" + key
}

func (s *FileStorage) Fetch(key string) (*Object, error) {
	key = s.BuildKey(key)

	_, err := os.Stat(key)
	if err != nil {
		return nil, ErrImageNotFound
	}

	fp, err := os.Open(key)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}
	defer fp.Close()

	logger.WithFields(logrus.Fields{
		"key": key,
	}).Debug("found object from file")

	object := &Object{}
	object.Body, err = ioutil.ReadAll(fp)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	// Ignore metadata file parse error
	metadataFp, err := os.Open(key + ".metadata")
	if err == nil {
		decorder := json.NewDecoder(metadataFp)
		err = decorder.Decode(&object.Metadata)
		if err != nil {
			logger.ErrorDebug(err)
		}
	} else {
		logger.ErrorDebug(err)
	}
	defer metadataFp.Close()

	return object, nil
}

func (s *FileStorage) PutFromBlob(key string, image []byte, contentType string, metadata map[string]string) error {
	key = s.BuildKey(key)

	directory := filepath.Dir(key)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return logger.ErrorDebug(err)
	}

	ioutil.WriteFile(key, image, os.ModePerm)

	if err != nil {
		return logger.ErrorDebug(err)
	}

	metadata["Content-Type"] = contentType

	j, err := json.Marshal(metadata)
	if err != nil {
		return logger.ErrorDebug(err)
	}
	ioutil.WriteFile(key+".metadata", j, os.ModePerm)

	logger.WithFields(logrus.Fields{
		"directory": directory,
		"key":       key,
	}).Debug("put to file")

	return nil
}

func (s *FileStorage) Put(key string, imageFile io.ReadSeeker, contentType string, metadata map[string]string) error {
	image, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return logger.ErrorDebug(err)
	}
	return s.PutFromBlob(key, image, contentType, metadata)
}

func (s *FileStorage) List(key string) ([]StorageItem, error) {
	path := s.BuildKey(key)

	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	items := make([]StorageItem, 0)
	for _, info := range fileInfos {
		logger.WithFields(logrus.Fields{
			"path": path,
			"name": info.Name(),
			"key":  key,
		}).Debug("found object")
		item := FileStorageItem{Name: key + "/" + info.Name()}
		items = append(items, &item)
	}

	return items, nil
}

func (s *FileStorage) Move(from string, to string) error {
	fromKey := s.BuildKey(from)
	toKey := s.BuildKey(to)

	logger.WithFields(logrus.Fields{
		"from": fromKey,
		"to":   toKey,
	}).Debug("move file object start")

	directory := filepath.Dir(fromKey)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return logger.ErrorDebug(err)
	}

	directory = filepath.Dir(toKey)
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		return logger.ErrorDebug(err)
	}

	err = os.Rename(fromKey, toKey)
	if err != nil {
		return logger.ErrorDebug(err)
	}

	return nil
}

func (s *FileStorageItem) IsValid() bool {
	if len(s.Extension()) == 0 {
		return false
	}

	if len(s.ImageSize()) == 0 {
		return false
	}

	return true
}

func (s *FileStorageItem) Key() string {
	return s.Name
}

func (s *FileStorageItem) Filename() string {
	path := strings.Split(s.Key(), "/")
	return path[len(path)-1]
}

func (s *FileStorageItem) Extension() string {
	path := strings.Split(s.Name, ".")
	return path[len(path)-1]
}

// KeyFormat: :image_type/:id/:id.:size.:format
func (s *FileStorageItem) ImageSize() string {
	path := strings.Split(s.Name, ".")
	return path[len(path)-2]
}
