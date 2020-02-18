package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

type OrganizationsConnector interface {
	CreateAccount(*organizations.CreateAccountInput) (*organizations.CreateAccountOutput, error)
}

type AccountCreationRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateAccount(oc OrganizationsConnector, ar AccountCreationRequest) (string, error) {
	out, err := oc.CreateAccount(&organizations.CreateAccountInput{
		AccountName: aws.String(ar.Name),
		Email:       aws.String(ar.Email),
	})
	if err != nil {
		return "", err
	}
	return *out.CreateAccountStatus.Id, nil
}

type CreatedAccount struct {
	StatusID string `json:"status_id"`
}

func createAccountHandler(sess *session.Session, req *events.APIGatewayProxyRequest) (CreatedAccount, error) {
	var ar AccountCreationRequest
	var out CreatedAccount

	if err := json.Unmarshal([]byte(req.Body), &ar); err != nil {
		return out, err
	}

	orgs := organizations.New(sess)
	id, err := CreateAccount(orgs, ar)
	if err != nil {
		return out, err
	}
	out.StatusID = id
	return out, nil
}
