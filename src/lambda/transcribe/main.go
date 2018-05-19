package main

import (
	"context"
	"encoding/json"
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
	//Marshal the eventinfo
	data, _ := json.Marshal(s3Event)
	streventinfo := string(data)

	// stdout and stderr are sent to AWS CloudWatch Logs
	fmt.Printf("S3 Event : %s\n", streventinfo)
	// interate through each record entry in the event data
	for _, record := range s3Event.Records {
		s3 := record.S3
		key := s3.Object.Key
		//obtain the position that represents the file suffix, including the period.
		index := (len(key)) - 4
		// create a substring that represents from index position to to the end of the string
		substring := key[index:]
		fmt.Printf("The file suffix is : %s\n", substring)
		// terminate script is the file is not an mp4 file
		if strings.ToUpper(substring) != ".MP4" {
			fmt.Printf("The object %s is not an mp4 file, exiting", s3.Object.Key)
			return
		}

		// open a new session
		sess, _ := session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String("eu-west-1")},
			Profile: "development",
		})

		log.Printf("Opening Transcribe session\n")
		transcriber := transcribeservice.New(sess)

		// exit if unable to create a Transcribe session
		if transcriber == nil {
			log.Printf("Unable to create Transcribe session\n")
			return
		} else {
			log.Printf("Transcribe session successfully created\n")
		}

		// create a random id for the jobname
		jobname := misc.GUID()
		mediafileuri := fmt.Sprintf("https://s3-eu-west-1.amazonaws.com/%s/%s", s3.Bucket.Name, s3.Object.Key)
		log.Printf("Job name :  %s\nMediaFileUri : %s\n", jobname, mediafileuri)

		mediaformat := "mp4"
		languagecode := "en-US"

		log.Printf("Creating transcription job\n")

		var StrucMedia transcribeservice.Media
		StrucMedia.MediaFileUri = &mediafileuri

		transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
			TranscriptionJobName: &jobname,
			Media:                &StrucMedia,
			MediaFormat:          &mediaformat,
			LanguageCode:         &languagecode,
		})
	}
	log.Printf("Complete")

}

func main() {
	lambda.Start(Handler)
}
