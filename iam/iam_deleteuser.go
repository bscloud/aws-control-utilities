package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Usage:
// go run iam_deleteuser.go <username>
func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := iam.New(sess)

	_, err = svc.DeleteUser(&iam.DeleteUserInput{
		UserName: &os.Args[1],
	})

	// If the user does not exist than we will log an error.
	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		fmt.Printf("User %s does not exist\n", os.Args[1])
		return
	} else if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Printf("User %s has been deleted\n", os.Args[1])
}
