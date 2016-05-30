package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"net/http"
)

type S3ReadAuth struct {
	AbstractAuth
}

func (a S3ReadAuth) SuccessLogLevel() string {
	return a.GetLogLevel("auth.s3_read_auth.success_log_level", "Info")
}

func (a S3ReadAuth) FailureLogLevel() string {
	return a.GetLogLevel("auth.s3_read_auth.failure_log_level", "Warn")
}

func (a S3ReadAuth) Authenticate(r *http.Request, ps httprouter.Params) bool {
	accessKeyId := r.FormValue("access_key_id")
	secretAccessKey := r.FormValue("secret_access_key")
	sessionToken := r.FormValue("session_token")
	cred := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, sessionToken)
	client := s3.New(awsSession.New(), &aws.Config{Credentials: cred, Region: aws.String(os.Getenv("KINU_S3_REGION"))})
	params := &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("KINU_S3_BUCKET")),
		Key:    aws.String("kinu_access_check"),
	}
	_, err := client.GetObject(params)
	if reqerr, ok := err.(awserr.RequestFailure); ok && reqerr.StatusCode() == http.StatusNotFound {
		return true
	} else if err != nil {
		return false
	}
	return true
}

func (a S3ReadAuth) HandleAuthError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) bool {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return true
}

