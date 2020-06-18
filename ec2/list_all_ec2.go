package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	/*if len(os.Args) != 2 {
		exitErrorf("Profile name required\nUsage: %s profile_name", os.Args[0])
	}

	profile := os.Args[1]*/

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "PROFILE",

		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error)
		os.Exit(1)
	}

	ec2Svc := ec2.New(sess)

	// call to get detailed information on each instance
	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Success", result)
	}

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
