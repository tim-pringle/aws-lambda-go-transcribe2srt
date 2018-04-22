package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

var (
	// ErrJobnameNotProvided is thrown when a name is not provided
	ErrJobnameNotProvided = errors.New("No job number was provided in the HTTP body")
	ErrTranscribeRunning  = errors.New("Job is still running")
	ErrTranscribeFailure  = errors.New("There was a problem running the job")
	ErrGeneral            = errors.New("An unknown error occurred")
	jobname               string
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	sess, _ := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-west-1")},
		Profile: "development",
	})

	if len(request.Body) < 1 {
		//return events.APIGatewayProxyResponse{}, ErrJobnameNotProvided
		log.Printf("No content in body, using default job number")
		jobname = "01524-38742-85260-18023"
	} else {
		log.Printf("Job number received %s", jobname)
		jobname = request.Body
	}

	log.Printf("Creating new session")
	transcriber := transcribeservice.New(sess)

	log.Printf("Getting transcription job")
	transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
		TranscriptionJobName: &jobname,
	})

	if err != nil {
		log.Printf("Unable to get transcription job %s", jobname)
		log.Printf(err.Error())
		return events.APIGatewayProxyResponse{}, ErrGeneral
	}
	strStatus := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
	log.Printf("Job status is %s", strStatus)

	var strFailureReason *string

	if strStatus == "FAILED" {
		strFailureReason = transcriptionjoboutput.TranscriptionJob.FailureReason
		return events.APIGatewayProxyResponse{
			Body:       *strFailureReason,
			StatusCode: 200,
		}, ErrTranscribeFailure
	}
	if strStatus == "IN_PROGRESS" {
		return events.APIGatewayProxyResponse{}, ErrTranscribeRunning
	}

	var uri *string
	uri = transcriptionjoboutput.TranscriptionJob.Transcript.TranscriptFileUri
	log.Printf("URI is %s", *uri)

	response, _ := http.Get(*uri)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	//str := string(body[:])

	//If there's an error, print the error
	if err != nil {
		fmt.Println(err)
	}

	// initialize our variable to hold the json
	var awstranscript Awstranscript

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'awstranscript' which we defined above
	json.Unmarshal(body, &awstranscript)

	var transcription []Item
	transcription = awstranscript.Results.Items

	var index, sequence int = 0, 0
	var srtinfo, subdetail, subtitle, sttime, classification, text, entime string
	var strlen int
	var firstrow bool

	for index = 0; index < len(transcription); {
		//Variable initiation for length of subtitle text, sequence number if its the first line and the subtitle text

		sequence++
		firstrow = true
		subtitle = ""

		//Grab the start time, convert it to a number, then convert the number an SRT valid time string
		sttime = transcription[index].Starttime
		fsttime, err := strconv.ParseFloat(sttime, 64)
		if err != nil {
			fmt.Println(err)
		}
		sttime = getsrttime(fsttime)

		/*Repeat this until we have either reached the last item in results
		#or the length of the lines we are reading is greater than 64 characters */

		for strlen = 0; (strlen < 64) && (index < len(transcription)); {
			text = transcription[index].Alternatives[0].Content
			strlen += len(text)

			switch classification {

			case "punctuation":
				if len(subtitle) > 0 {
					runes := []rune(subtitle)
					subtitle = string(runes[1 : len(subtitle)-1])
				} else {
					subtitle += text
				}
			default:
				subtitle += (text + " ")
			}

			//If the length of the current string is greater than 32 and this
			//is the first line of the sequence, then add a return character to it

			if (strlen > 32) && (firstrow == true) {
				subtitle += "\n"
				firstrow = false
			}

			/*If the last character is a '.', then we need to set
			the end time attribute to the previous indexes one
			since punctuation characters to not have a time stamp*/

			if classification == "punctuation" {
				entime = transcription[index-1].Endtime
			} else {
				entime = transcription[index].Endtime
			}

			fsttime, err = strconv.ParseFloat(entime, 64)
			entime = getsrttime(fsttime)

			index++
		}
		//Setup the string that is refers to these two
		//lines in SRT format

		subdetail = fmt.Sprintf("\n%d\n%s --> %s\n%s\n", sequence, sttime, entime, subtitle)

		//Append this to the existing string
		srtinfo += subdetail

	}

	log.Printf(srtinfo)

	return events.APIGatewayProxyResponse{
		Body:       srtinfo,
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
