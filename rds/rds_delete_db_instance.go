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
		exitErrorf("instance name required\nUsage: %s instance_name", os.Args[0])
	}

	sqli := os.Args[1]

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

	input := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(sqli),
		SkipFinalSnapshot:    aws.Bool(true),
	}

	result, err := svc.DeleteDBInstance(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				fmt.Println(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
			case rds.ErrCodeInvalidDBInstanceStateFault:
				fmt.Println(rds.ErrCodeInvalidDBInstanceStateFault, aerr.Error())
			case rds.ErrCodeDBSnapshotAlreadyExistsFault:
				fmt.Println(rds.ErrCodeDBSnapshotAlreadyExistsFault, aerr.Error())
			case rds.ErrCodeSnapshotQuotaExceededFault:
				fmt.Println(rds.ErrCodeSnapshotQuotaExceededFault, aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				fmt.Println(rds.ErrCodeInvalidDBClusterStateFault, aerr.Error())
			case rds.ErrCodeDBInstanceAutomatedBackupQuotaExceededFault:
				fmt.Println(rds.ErrCodeDBInstanceAutomatedBackupQuotaExceededFault, aerr.Error())
			default:
				fmt.Println(aerr.Error)
			}
		} else {
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
