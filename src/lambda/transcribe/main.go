package main

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	jobname := GUID()
	log.Printf("Job name created %s\n", jobname)

	joburi := "https://s3-eu-west-1.amazonaws.com/tim-training-thing/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.mp4"
	mediaformat := "mp4"
	languagecode := "en-US"

	sess, _ := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-west-1")},
		Profile: "development",
	})

	log.Printf("Opening Transcribe session %s\n", request.RequestContext.RequestID)
	transcriber := transcribeservice.New(sess)
	if transcriber == nil {
		log.Printf("Unable to create Transcribe session %s\n", request.RequestContext.RequestID)
	} else {
		log.Printf("Transcribe session successfully created %s\n", request.RequestContext.RequestID)
	}

	log.Printf("Creating transcription job %s\n", request.RequestContext.RequestID)

	var StrucMedia transcribeservice.Media
	StrucMedia.MediaFileUri = &joburi

	transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
		TranscriptionJobName: &jobname,
		Media:                &StrucMedia,
		MediaFormat:          &mediaformat,
		LanguageCode:         &languagecode,
	})

	transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
		TranscriptionJobName: &jobname,
	})

	if err != nil {
		log.Printf("Unable to get job output %s\n", request.RequestContext.RequestID)
	} else {
		log.Printf("Retrieved job output %s\n", request.RequestContext.RequestID)
	}

	strStatus := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
	strJobname := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobName)
	var strFailureReason *string

	if strStatus == "FAILED" {
		strFailureReason = transcriptionjoboutput.TranscriptionJob.FailureReason
		var ErrTransctibeFailure = errors.New(*strFailureReason)
		return events.APIGatewayProxyResponse{}, ErrTransctibeFailure
	}
	//Success is denoted by the function not already having exited via error checks
	return events.APIGatewayProxyResponse{
		Body:       strJobname,
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
