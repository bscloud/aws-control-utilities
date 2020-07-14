package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	if len(os.Args) != 2 {
		exitErrorf("rds cluster name required\nUser: %s cluster_name", os.Args[0])
	}

	cn := os.Args[1]

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "kroton-captacao-dev",

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

	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(cn),
	}

	result, err := svc.DescribeDBClusters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				fmt.Println(rds.ErrCodeDBClusterNotFoundFault, aerr.Error())
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
