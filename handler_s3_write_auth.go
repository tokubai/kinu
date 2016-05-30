package main

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"net/http"
)

type S3WriteAuth struct {
	AbstractAuth
}

func (a S3WriteAuth) SuccessLogLevel() string {
	return a.GetLogLevel("auth.s3_write_auth.success_log_level", "Info")
}

func (a S3WriteAuth) FailureLogLevel() string {
	return a.GetLogLevel("auth.s3_write_auth.failure_log_level", "Warn")
}

func (a S3WriteAuth) Authenticate(r *http.Request, ps httprouter.Params) bool {
	accessKeyId := r.FormValue("access_key_id")
	secretAccessKey := r.FormValue("secret_access_key")
	sessionToken := r.FormValue("session_token")
	cred := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, sessionToken)
	client := s3.New(awsSession.New(), &aws.Config{Credentials: cred, Region: aws.String(os.Getenv("KINU_S3_REGION"))})
	_, err := client.PutObject(&s3.PutObjectInput{
		Bucket:   aws.String(os.Getenv("KINU_S3_BUCKET")),
		Key:      aws.String("kinu_write_check"),
		Body:     bytes.NewReader([]byte("dummy")),
	})

	if err != nil {
		return false
	}
	return true
}

func (a S3WriteAuth) HandleAuthError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) bool {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return true
}

