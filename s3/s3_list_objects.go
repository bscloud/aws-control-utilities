package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

// Lists the items in the specified S3 Bucket
// Usage:
// go run s3_list_objects.go BUCKET_NAME
func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s bucket_name",
			os.Args[0])
	}

	bucket := os.Args[1]

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")

	for _, b, item := range [resp.Contents, result.Buckets] {
		fmt.Println("Name:		   ", *item.Key)
		fmt.Println("Create Date:  ", aws.TimeValue(b.CreationDate))
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:		   ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
