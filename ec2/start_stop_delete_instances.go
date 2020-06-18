package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Usage:
// go run main.go <state> <instance id>
//   * state can either be START or STOP
func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create new EC2 client
	svc := ec2.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	// Instances delete
	if os.Args[1] == "DELETE" {
		input := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.TerminateInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			// Let's now set dry run to be false. This will allow us to start the instances
			input.DryRun = aws.Bool(false)
			result, err = svc.TerminateInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Success", result)
			}
		} else { // This could be due to a lack of permissions
			fmt.Println("Error", err)
		}
	}

	// Turn monitoring on
	if os.Args[1] == "START" {
		// We set DryRun to true to check to see if the instance exists and we have the
		// necessary permissions to monitor the  instance
		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StartInstances(input)
		awsErr, ok := err.(awserr.Error)

		// If the error code is `DryRunOperation` it means we have the necessary
		// permissions to Start this instance
		if ok && awsErr.Code() == "DryRunOperation" {
			// Let's now set dry run to be false. This will allow us to start the instances
			input.DryRun = aws.Bool(false)
			result, err = svc.StartInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Success", result.StartingInstances)
			}
		} else { // This could be due to a lack of permissions
			fmt.Println("Error", err)
		}
	} else if os.Args[1] == "STOP" { // Turn instances off
		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.StopInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Success", result.StoppingInstances)
			}
		} else {
			fmt.Println("Error", err)
		}
	}
}
