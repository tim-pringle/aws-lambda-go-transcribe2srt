package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	var request AlexaRequest
	err := json.Unmarshal(data, &request)
	log.Printf("A request from Alexa has been received")
	if err != nil {
		log.Printf("Problem encountered unmarshalling request to Alex Response struct")
	}

	var response AlexaResponse
	var arrdialog []DialogDelegate
	var dialog DialogDelegate

	//Depending on the dialog state there needs to be a different action.
	//STARTED --> IN_PROGRESS --> xxxx --> COMPLETED
	if request.Request.DialogState == "STARTED" {
		log.Printf("DialogueState is STARTED, requesting delegation")
		dialog.Type = "Dialog.Delegate"
		dialog.UpdatedIntent = nil
		arrdialog = append(arrdialog, dialog)
		response.Version = "1.0"
		response.Response.OutputSpeech = nil
		response.Response.Card = nil
		response.Response.Reprompt = nil
		response.Response.Directives = &arrdialog
		strResponse, _ := json.Marshal(response)
		log.Printf("-----------------DELEGATION REQUEST-------------------")
		log.Printf("%+v", string(strResponse))
		log.Printf("-----------------DELEGATION REQUEST-------------------")
		return response, nil

	} else if request.Request.DialogState == "IN_PROGRESS" {
		log.Printf("DialogueState is IN_PROGRESS, parsing job data")
		jbname := request.Request.Intent.Slots.Jobnumber.Value
		//You must use a custom slot format - Alexa does not accept numbers of characters in length
		if len(jbname) < 20 {
			var os OutputSpeech

			os.Text = fmt.Sprintf("The suppled job number of %s is not valid", jbname)
			os.Type = "PlainText"

			response.Version = "1.0"
			response.Response.OutputSpeech = &os
			return response, nil

		}
		log.Printf("Received: %s", jbname)
		jobname = fmt.Sprintf("%s-%s-%s-%s", jbname[0:5], jbname[5:10], jbname[10:15], jbname[15:20])
		if err != nil {
			log.Printf("%s", err.Error())
			return err.Error(), nil
		}

		//log.Printf("%s", jobname)
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

		strResponse, _ := json.Marshal(response)
		log.Printf("-----------------IN PROCESS REQUEST-------------------")
		log.Printf("%+v", string(strResponse))
		log.Printf("-----------------IN PROCESS REQUEST-------------------")

	}
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
