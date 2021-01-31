package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/lambda"
)

//
type Data struct {
	Body string
}

type Request struct {
	Query string
}

type Response struct {
	Data *Data `json:"body"`
}

func HandleLambdaEvent(event Request) (Response, error) {
	d, err := parseQuery(event.Query)
	if err != nil {
		return Response{Data: d}, nil
	}
	return Response{}, errors.New("error occured")
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
