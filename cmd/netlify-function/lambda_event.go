package main

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

func handleLambdaEvent(handler http.HandlerFunc, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	r, err := http.NewRequest(event.HTTPMethod, event.Path, strings.NewReader(event.Body))
	if err != nil {
		return nil, errors.Wrap(err, "error converting lambda event to http request")
	}

	w := respWriter{
		header: make(http.Header),
	}

	handler(&w, r)

	return &events.APIGatewayProxyResponse{
		StatusCode: w.statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "X-Requested-With,Content-Type,Authorization,Auth-Token",
			"Access-Control-Allow-Methods": "GET,PUT,POST,DELETE,OPTIONS,PING",
		},
		Body: string(w.b),
	}, nil
}
