package main

import (
	"os"
)

type Config struct {
	Bucket   string
	Region   string
	BasePath string
}

var (
	config *Config
)

func init() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		panic("config load failed")
	}
}

func LoadConfig() (*Config, error) {
	bucket, _ := os.LookupEnv("KINU_S3_BUCKET")
	basePath, _ := os.LookupEnv("KINU_S3_BUCKET_BASE_PATH")
	region, _ := os.LookupEnv("KINU_S3_REGION")
	return &Config{Bucket: bucket, BasePath: basePath, Region: region}, nil
}
