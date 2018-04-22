package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {

	//Set the S3 uri prefix and create a unique guid to be used as the job name

	bucket := flag.String("bucket", "tim-training-thing", "The s3 bucket to upload to")
	filename := flag.String("filename", "", "The file to be uploaded")

	//Lets set default parameter values for the aws cmdlets

	sess, _ := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-west-1")},
		Profile: "development",
	})

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(*filename)
	key := filepath.Base(file.Name())

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		fmt.Println(err)
	}
}
