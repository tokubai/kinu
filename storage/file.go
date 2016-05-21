package storage

import (
	"github.com/Sirupsen/logrus"
	"github.com/TakatoshiMaeda/kinu/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func (s *FileStorage) Fetch(key string) ([]byte, error) {
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

	image, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	return image, nil
}

func (s *FileStorage) PutFromBlob(key string, image []byte) error {
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

	logger.WithFields(logrus.Fields{
		"directory": directory,
		"key":       key,
	}).Debug("put to file")

	return nil
}

func (s *FileStorage) Put(key string, imageFile io.ReadSeeker) error {
	image, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return logger.ErrorDebug(err)
	}
	return s.PutFromBlob(key, image)
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

func (s *FileStorageItem) Extension() string {
	path := strings.Split(s.Name, ".")
	return path[len(path)-1]
}

// KeyFormat: :image_type/:id/:id.:size.:format or :image_type/:id/:id.:format
func (s *FileStorageItem) ImageSize() string {
	if sizeHasImageFileNameRegexp.MatchString(s.Name) {
		path := strings.Split(s.Name, ".")
		return path[len(path)-2]
	} else {
		return "1000"
	}
}
