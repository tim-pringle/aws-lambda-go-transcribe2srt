package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tim-pringle/transcribe2srt"
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

	//Lex Request
	var lexrq events.LexEvent
	json.Unmarshal(data, &lexrq)

	jobname = lexrq.InputTranscript
	subtitles, converterror := transcribe2srt.Convert(jobname)
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

}

func main() {
	lambda.Start(Handler)
}
