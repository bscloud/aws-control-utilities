package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ListBuckets() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "kroton-analytics-stg",

		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		fmt.Println("Go erro creating new session")
		fmt.Println(err.Error)
		os.Exit(1)
	}

	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", &bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name: 					", *item.Key)
		fmt.Println("Last Modified: ", *item.LastModified)
		fmt.Println("Size: 				: ", *item.Size)
		fmt.Println("Storage Class: ", *item.Size)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "Items in bucket", bucket)
	fmt.Println("")
}

func main() {
	if len(os.Args) != 2 {
		exitErrorf("bucket name required\nUsage: %s buckt_name", os.Args[0])
	}

	bucket := os.Args[1]

	ListBuckets(bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
