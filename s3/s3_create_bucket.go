package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name missing!\nUsege: % bucket_name", os.Args[0])
	}

	bucket := os.Args[1]

	/*sess. err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")
	})*/

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "kroton-analytics-stg",

		Config: aws.Config{
			Region: aws.String("us-west-2"),
		},

		//SharedConfigState: SharedConfigEnable,
	})

	svc := s3.New(sess)

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	fmt.Printf("Bucket: %q successfully created\n", bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
