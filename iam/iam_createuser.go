package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Usage:
// go run iam_createuser.go <username>
func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := iam.New(sess)

	_, err := svc.GetUser(&iam.GetUserInput{
		UserName: &os.Args[1],
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		result, err := svc.CreateUser(&iam.CreateUserInput{
			UserName: &os.Args[1]
		})

		if err != nil {
			fmt.Println("CreateUser Error", err)
			return
		}

		fmt.Println("Success", result)
	} else {
		fmt.Println("Getuser Error", err)
	}
}
