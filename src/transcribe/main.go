package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

//GUID - generates a unique identifier
func GUID() (guid string) {
	ad, _ := time.Parse("02-01-2006", "01-01-1970")

	timesince := time.Since(ad).Nanoseconds()
	strsince := strconv.FormatInt(timesince, 10)
	guid = fmt.Sprintf("0" + strsince[0:4] + "-" + strsince[4:9] + "-" + strsince[9:14] + "-" + strsince[14:19])
	return
}

// Handler is the Lambda function handler
// It uses an S3 event source, with the Lambda function being trigged
// when a CreateObject event occurs on an S3 bucket that has Events configured
func Handler(ctx context.Context, s3Event events.S3Event) {
	//Marshal the eventinfo
	data, _ := json.Marshal(s3Event)
	//Now convert to a string and output
	//Cloudwatch picks up the json and formats it nicely for us. :)
	streventinfo := string(data)

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("S3 Event : %s\n", streventinfo)
	// interate through each record entry in the event data
	for _, record := range s3Event.Records {
		s3 := record.S3

		log.Printf("Object : %s\n", s3.Object.Key)

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
		jobname := GUID()
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
