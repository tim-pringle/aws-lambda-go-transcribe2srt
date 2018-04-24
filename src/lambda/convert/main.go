package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	jobname string
)

type LexResponse struct {
	SessionAttributes struct {
	} `json:"sessionAttributes"`
	DialogAction struct {
		Type             string `json:"type"`
		FulfillmentState string `json:"fulfillmentState"`
		Message          struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"message"`
	} `json:"dialogAction"`
}

func Handler(ctx context.Context, eventinfo interface{}) (interface{}, error) {

	//Marshal the eventinfo
	data, _ := json.Marshal(eventinfo)
	streventinfo := string(data)

	//Lets try and cast this into a WebProxy Request

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
