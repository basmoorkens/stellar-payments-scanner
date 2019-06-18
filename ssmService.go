package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/sirupsen/logrus"
)

func getLastScannedHorizonCursor() string {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "tradedirect"),
	})
	if err != nil {
		panic(err)
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	keyname := "stellar-last-read-cursor"
	withDecryption := false
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	value := *param.Parameter.Value
	return value
}

func saveCursor(cursorVal string) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "tradedirect"),
	})
	if err != nil {
		panic(err)
	}
	logrus.Info("Saving new cursor value " + cursorVal)
	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	keyname := "stellar-last-read-cursor"
	paramType := "String"
	resp, err := ssmsvc.PutParameter(&ssm.PutParameterInput{
		Name:      &keyname,
		Value:     &cursorVal,
		Overwrite: aws.Bool(true),
		Type:      &paramType,
	})
	if err != nil {
		panic(err)
	}
	fmt.Print(resp.Version)
}
