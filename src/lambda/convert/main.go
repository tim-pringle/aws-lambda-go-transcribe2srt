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

func Handler(ctx context.Context, eventinfo interface{}) (interface{}, error) {

	//Marshal the eventinfo
	data, _ := json.Marshal(eventinfo)
	streventinfo := string(data)

	data, _ = json.Marshal(ctx)
	strcontextinfo := string(data)

	//Lets log the request for reference
	log.Printf("-----------------CONTEXT INFO-----------------")
	log.Printf(streventinfo)
	log.Printf("-----------------CONTEXT INFO-----------------")
	log.Printf("")
	log.Printf("-----------------EVENT INFO-------------------")
	log.Printf(strcontextinfo)
	log.Printf("-----------------EVENT INFO-------------------")
	log.Printf("")

	//Lets have a look at the request

	//Alex Request
	if (strings.Contains(streventinfo, "LaunchRequest")) && (strings.Contains(streventinfo, "amazonalexa")) {
		jobname = "01524-63806-28098-90622"
		subtitles, converterror := Convert(jobname)

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
