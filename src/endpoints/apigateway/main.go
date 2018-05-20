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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//Marshal the eventinfo for nice formatting in the logs
	data, _ := json.Marshal(request)
	streventinfo := string(data)

	//Lets log the request for reference
	log.Printf("-----------------EVENT INFO-------------------")
	log.Printf("%+v", streventinfo)
	log.Printf("-----------------EVENT INFO-------------------")

	// API Gateway Response
	var response events.APIGatewayProxyResponse
	if len(request.Body) < 1 {
		response.Body = "Empty request body"
		response.StatusCode = 400
		return response, nil
	}

	json.Unmarshal(data, &request)

	jobname = request.Body
	subtitles, converterror := transcribe2srt.Convert(jobname)

	if converterror != nil {
		response.Body = "Server error"
		response.StatusCode = 400
	} else {
		response.Body = subtitles
		response.StatusCode = 200
	}
	return response, nil

}

func main() {
	lambda.Start(handler)
}
