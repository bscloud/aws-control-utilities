package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	if len(os.Args) != 2 {
		exitErrorf("bucket name required\nUsage: %s buckt_name", os.Args[0])
	}

	bucket := os.Args[1]

	sess, _ := session.NewSessionWithOptions(session.Options{
		Profile: "kroton-analytics-stg",

		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})

	svc := s3.New(sess)

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		exitErrorf("Unable to delete objects from %q, %v", bucket, err)
	}

	fmt.Printf("Deleted object(s) from bucket: %s", bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
