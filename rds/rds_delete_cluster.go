package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	if len(os.Args) != 2 {
		exitErrorf("rds cluster name required\nUser: %s cluster_name", os.Args[0])
	}

	cn := os.Args[1]

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

	input := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(cn),
		SkipFinalSnapshot:   aws.Bool(true),
	}

	result, err := svc.DeleteDBCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				fmt.Println(rds.ErrCodeDBClusterNotFoundFault, aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				fmt.Println(rds.ErrCodeInvalidDBClusterStateFault, aerr.Error())
			case rds.ErrCodeDBClusterSnapshotAlreadyExistsFault:
				fmt.Println(rds.ErrCodeDBClusterSnapshotAlreadyExistsFault, aerr.Error())
			case rds.ErrCodeSnapshotQuotaExceededFault:
				fmt.Println(rds.ErrCodeSnapshotQuotaExceededFault, aerr.Error())
			case rds.ErrCodeInvalidDBClusterSnapshotStateFault:
				fmt.Println(rds.ErrCodeInvalidDBClusterSnapshotStateFault, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
