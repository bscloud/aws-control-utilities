package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "kroton-captacao-stg",

		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error)
		os.Exit(1)
	}

	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(nil)
	if err != nil {
		exitErrorf("Unable to list instances, %v", err)
	}

	fmt.Println("Instances:")

	for _, d := range result.DBInstances {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(d.DBInstanceIdentifier), aws.TimeValue(d.InstanceCreateTime))
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
