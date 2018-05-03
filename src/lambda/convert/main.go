package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/tim-pringle/go-aws/transcribe"
)

var (
	jobname string
)

func Handler(ctx context.Context, eventinfo interface{}) (interface{}, error) {

	//Marshal the eventinfo
	data, _ := json.Marshal(eventinfo)
	streventinfo := string(data)

	//Lets log the request for reference
	log.Printf("-----------------EVENT INFO-------------------")
	log.Printf("%+v", streventinfo)
	log.Printf("-----------------EVENT INFO-------------------")

	//Lets have a look at the request
	//Alex Request
	if (strings.Contains(streventinfo, "LaunchRequest")) || (strings.Contains(streventinfo, "IntentRequest")) || (strings.Contains(streventinfo, "SessionEndedRequest")) {
		var request AlexaRequest
		err := json.Unmarshal(data, &request)
		log.Printf("Alexa request received")
		if err != nil {
			log.Printf("Problem encountered unmarshalling request to Alex Response struct")
		}

		var response AlexaResponse
		var arrdialog []DialogDelegate
		var dialog DialogDelegate

		//Depending on the dialog state there needs to be a different action.
		//STARTED --> IN_PROGRESS --> xxxx --> COMPLETED
		if request.Request.DialogState == "STARTED" {
			dialog.Type = "Dialog.Delegate"
			dialog.UpdatedIntent = nil
			arrdialog = append(arrdialog, dialog)
			response.Version = "1.0"
			response.Response.OutputSpeech = nil
			response.Response.Card = nil
			response.Response.Reprompt = nil
			response.Response.Directives = &arrdialog
		} else if request.Request.DialogState == "IN_PROGRESS" {
			jobname = "01524-63806-28098-90622"
			//jobname := request.Request.Intent.Slots.Jobnumber.Value
			log.Printf("Job number received : %+v", jobname)

			log.Printf("Processing conversion request")

			sess, _ := session.NewSessionWithOptions(session.Options{
				Config:  aws.Config{Region: aws.String("eu-west-1")},
				Profile: "development",
			})

			log.Printf("Job name : %s", jobname)
			log.Printf("Creating new session")
			transcriber := transcribeservice.New(sess)

			log.Printf("Getting transcription job")
			transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
				TranscriptionJobName: &jobname,
			})

			if err != nil {
				ErrMsg := errors.New(err.Error())
				log.Printf("%s", err.Error())
				return "", ErrMsg
			}

			strStatus := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
			log.Printf("Job status is %s", strStatus)

			if strStatus == "FAILED" {
				return "", ErrTranscribeFailure
			}
			if strStatus == "IN_PROGRESS" {
				return "", ErrTranscribeRunning
			}

			var uri *string
			uri = transcriptionjoboutput.TranscriptionJob.Transcript.TranscriptFileUri
			log.Printf("URI is %s", *uri)

			resp, _ := http.Get(*uri)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			//str := string(body[:])

			//If there's an error, print the error
			if err != nil {
				fmt.Println(err)
			}

			log.Printf(string(body))

			// initialize our variable to hold the json
			var awstranscript transcribe.Awstranscript
			var index = 0
			// we unmarshal our byteArray which contains our
			// jsonFile's content into 'awstranscript' which we defined above
			json.Unmarshal(body, &awstranscript)
			transcription := awstranscript.Results.Transcripts
			var text = ""
			for index = 0; index < len(transcription); {
				text = transcription[index].Transcript
				index++
			}
			log.Printf("The text will be : %s", text)

			var os OutputSpeech

			os.Text = text
			os.Type = "PlainText"

			response.Version = "1.0"
			response.Response.OutputSpeech = &os
			//response.Response.Card = nil
			//response.Response.Reprompt = nil
			log.Printf("%+v", response)
			//response.Response.Directives = nil

		}

		strResponse, _ := json.Marshal(response)
		log.Printf("-----------------LAMBDA RESPONSE-------------------")
		log.Printf("%+v", string(strResponse))
		log.Printf("-----------------LAMBDA REPSONSE-------------------")
		return response, nil
		//jobname = "01524-63806-28098-90622"
		//subtitles, converterror := Convert(jobname)

	}
	// API Gateway Request
	if (strings.Contains(streventinfo, "httpMethod")) && (strings.Contains(streventinfo, "headers")) {
		var request events.APIGatewayProxyRequest
		err := json.Unmarshal(data, &request)
		if len(request.Body) < 1 {
			log.Printf("No content in body")
			return "", err
		}
		jobname = request.Body
		subtitles, converterror := Convert(jobname)
		var response events.APIGatewayProxyResponse
		if converterror != nil {
			response.Body = "Server error"
			response.StatusCode = 400
		} else {
			response.Body = subtitles
			response.StatusCode = 200
		}

		return response, nil
		//Lex Request
	} else if (strings.Contains(streventinfo, "currentIntent")) && (strings.Contains(streventinfo, "userId")) {
		var lexrq events.LexEvent
		err := json.Unmarshal(data, &lexrq)
		if err != nil {

		}
		jobname = lexrq.InputTranscript
		subtitles, converterror := Convert(jobname)
		var response LexResponse

		if converterror != nil {
			response.DialogAction.Type = "Close"
			response.DialogAction.FulfillmentState = "Fulfilled"
			response.DialogAction.Message.ContentType = "PlainText"
			response.DialogAction.Message.Content = "Error!!!!!"
			return response, converterror
		} else {
			response.DialogAction.Type = "Close"
			response.DialogAction.FulfillmentState = "Fulfilled"
			response.DialogAction.Message.ContentType = "PlainText"
			response.DialogAction.Message.Content = subtitles
			return response, nil
		}

	} else {
		unsupported := errors.New("Unsupported service")
		return nil, unsupported
	}
}

func main() {
	lambda.Start(Handler)
}
