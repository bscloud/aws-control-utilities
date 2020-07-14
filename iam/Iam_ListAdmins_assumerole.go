package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

func UserPolicyHasAdmin(user *iam.UserDetail, admin string) bool {
	for _, policy := range user.UserPolicyList {
		if *policy.PolicyName == admin {
			return true
		}
	}

	return false
}

func AttachedUserPolicyHasAdmin(user *iam.UserDetail, admin string) bool {
	for _, policy := range user.AttachedManagedPolicies {
		if *policy.PolicyName == admin {
			return true
		}
	}

	return false
}

func GroupPolicyHasAdmin(svc *iam.IAM, group *iam.Group, admin string) bool {
	input := &iam.ListGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	result, err := svc.ListGroupPolicies(input)
	if err != nil {
		fmt.Println("Got error calling ListGroupPolicies fro group", group.GroupName)
	}

	// Wade through policies
	for _, policyName := range result.PolicyNames {
		if *policyName == admin {
			return true
		}
	}

	return false
}

func AttachedGroupPolicyHasAdmin(svc *iam.IAM, group *iam.Group, admin string) bool {
	input := &iam.ListAttachedGroupPoliciesInput{GroupName: group.GroupName}
	result, err := svc.ListAttachedGroupPolicies(input)
	if err != nil {
		fmt.Println("Got error gettin attached grou policies:")
		fmt.Println(err.Error)
		os.Exit(1)
	}

	for _, policy := range result.AttachedPolicies {
		if *policy.PolicyName == admin {
			return true
		}
	}

	return false
}

func UsersGroupsHaveAdmin(svc *iam.IAM, user *iam.UserDetail, admin string) bool {
	input := &iam.ListGroupsForUserInput{UserName: user.UserName}
	result, err := svc.ListGroupsForUser(input)
	if err != nil {
		fmt.Println("Got error getting groups for user:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, group := range result.Groups {
		groupPolicyHasAdmin := GroupPolicyHasAdmin(svc, group, admin)

		if groupPolicyHasAdmin {
			return true
		}

		attachedGroupPolicyHasAdmin := AttachedGroupPolicyHasAdmin(svc, group, admin)

		if attachedGroupPolicyHasAdmin {
			return true
		}
	}

	return false
}

func IsUserAdmin(svc *iam.IAM, user *iam.UserDetail, admin string) bool {
	// Check policy, attached policy, and groups (policy and attached policy)
	policyHasAdmin := UserPolicyHasAdmin(user, admin)
	if policyHasAdmin {
		return true
	}

	attachedPolicyHasAdmin := AttachedUserPolicyHasAdmin(user, admin)
	if attachedPolicyHasAdmin {
		return true
	}

	userGroupsHaveAdmin := UsersGroupsHaveAdmin(svc, user, admin)
	if userGroupsHaveAdmin {
		return true
	}

	return false
}

func main() {
	sess, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println("Got error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	stsSvc := sts.New(sess)
	stsCreadProvider := stscreds.NewAssumeRoleProvider(stsSvc, "arn:aws:iam::725186902151:role/RolePipeline")

	sess.Credentials = aws.NewCredentials(stsCreadProvider)

	svc := iam.New(sess)

	numUsers := 0
	numAdmins := 0

	user := "User"
	input := &iam.GetAccountAuthorizationDetailsInput{Filter: []*string{&user}}
	resp, err := svc.GetAccountAuthorizationDetails(input)
	if err != nil {
		fmt.Println("Got error getting account details")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// The policy name that indicate administrator acess
	adminName := "AdministratorAccess"

	// Wade through resulting users
	for _, user := range resp.UserDetailList {
		numUsers++

		isAdmin := IsUserAdmin(svc, user, adminName)

		if isAdmin {
			fmt.Println(*user.UserName)
			numAdmins += 1
		}
	}

	// Are there more?
	for *resp.IsTruncated {
		input := &iam.GetAccountAuthorizationDetailsInput{Filter: []*string{&user}, Marker: resp.Marker}
		resp, err = svc.GetAccountAuthorizationDetails(input)
		if err != nil {
			fmt.Println("Got error getting account details")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// wede through resulting user
		for _, user := range resp.UserDetailList {
			numUsers += 1

			isAdmin := IsUserAdmin(svc, user, adminName)

			if isAdmin {
				fmt.Println(*user.UserName)
				numAdmins += 1
			}
		}
	}

	fmt.Println("")
	fmt.Println("Found", numAdmins, "Admin(s) out of", numUsers, "user(s).")
}
