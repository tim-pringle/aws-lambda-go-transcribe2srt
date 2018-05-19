package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/tim-pringle/go-misc/misc"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

// Handler is the Lambda function handler
// It uses an S3 event source, with the Lambda function being trigged
// when a CreateObject event occurs on an S3 bucket that has Events configured
func Handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3

		fmt.Printf("[%s - %s] Bucket = %s, Key = %s  URL DecodedKey = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key, s3.Object.URLDecodedKey)
		var key string = s3.Object.Key
		index := (len(key)) - 4
		substring := key[index:]
		if strings.ToUpper(substring) != ".MP4" {
			fmt.Printf("The object %s is not an mp4 file", s3.Object.Key)
			return
		}
		// stdout and stderr are sent to AWS CloudWatch Logs

		jobname := misc.GUID()
		log.Printf("Job name created %s\n", jobname)

		joburi := fmt.Sprintf("https://s3-eu-west-1.amazonaws.com/%s/%s", s3.Bucket.Name, s3.Object.Key)
		log.Printf("Job uri : %s\n", joburi)

		mediaformat := "mp4"
		languagecode := "en-US"

		sess, _ := session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String("eu-west-1")},
			Profile: "development",
		})

		log.Printf("Opening Transcribe session\n")
		transcriber := transcribeservice.New(sess)
		if transcriber == nil {
			log.Printf("Unable to create Transcribe session\n")
		} else {
			log.Printf("Transcribe session successfully created\n")
		}

		log.Printf("Creating transcription job\n")

		var StrucMedia transcribeservice.Media
		StrucMedia.MediaFileUri = &joburi

		transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
			TranscriptionJobName: &jobname,
			Media:                &StrucMedia,
			MediaFormat:          &mediaformat,
			LanguageCode:         &languagecode,
		})
	}
}

func main() {
	lambda.Start(Handler)
}
