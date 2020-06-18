package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	// Return a list of key pairs
	result, err := svc.DescribeKeyPairs(nil)
	if err != nil {
		exitErrorf("Unable to get key pairs, %v", err)
	}

	fmt.Println("Key Pairs:")
	for _, pair := range result.KeyPairs {
		fmt.Printf("$s: %s\n", *pair.KeyName, *pair.KeyFingerprint)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
