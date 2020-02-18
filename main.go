package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"os"
)

type Action string

const (
	createAccount Action = "create-account"
)

type RequestType struct {
	Action string `json:"action"`
}

func requiredEnvVar(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic(fmt.Sprintf("Required environment variable %s not set", key))
}

func handler(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var rt RequestType

	roleArn := requiredEnvVar("ASSUME_ROLE")

	if err := json.Unmarshal([]byte(req.Body), &rt); err != nil {
		return validationErrorResponse(fmt.Sprintf("Error parsing the request: %s", err.Error()))
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{}))

	creds := stscreds.NewCredentials(sess, roleArn)
	sess = session.Must(session.NewSession(&aws.Config{Credentials: creds}))

	switch Action(rt.Action) {
	case createAccount:
		return successfulResponse(createAccountHandler(sess, req))
	}
	return validationErrorResponse("Unsupported action: " + rt.Action)
}

func successfulResponse(v interface{}, err error) (*events.APIGatewayProxyResponse, error) {
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	body, err := json.Marshal(v)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func validationErrorResponse(reason string) (*events.APIGatewayProxyResponse, error) {
	verr, _ := json.Marshal(ValidationError{Reason: reason})
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       string(verr),
	}, nil
}

type ValidationError struct {
	Reason string `json:"reason"`
}

func main() {
	lambda.Start(handler)
	// +create account

	// poll until account is created

	// move under labs OU

	// input: lab name, password

	// apply scp
	// grant root account's admin role access to the new account
	// create admin group
	// create user
	// send credentials
}
