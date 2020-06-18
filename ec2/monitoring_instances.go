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
//   * state can either be ON or OFF

func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create new EC2 client
	svc := ec2.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	// Turn monitoring on
	if os.Args[1] == "ON" {
		// We set DryRun to true to check to see if the instance exists and we have the
		// necessary permissions to monitor the instance.
		input := &ec2.MonitorInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.MonitorInstances(input)
		awsErr, ok := err.(awserr.Error)

		// If the error code is `DryRunOperation` it means we have the necessary
		// permissions to monitor this instance
		if ok && awsErr.Code() == "DryRunOperation" {
			// Let's now set dry run to be false. This will allow us to turn monitoring on
			input.DryRun = aws.Bool(false)
			result, err = svc.MonitorInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Sucess", result.InstanceMonitorings)
			}
		} else {
			// This could be due to a lack of permissions
			fmt.Println("Error", err)
		}
	} else if os.Args[1] == "OFF" { // Turn momitoring off
		input := &ec2.UnmonitorInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.UnmonitorInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.UnmonitorInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Sucess", result.InstanceMonitorings)
			}
		} else {
			fmt.Println("Error", err)
		}
	}
}
