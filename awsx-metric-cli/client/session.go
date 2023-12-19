package client

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/awssession"
	util "github.com/Appkube-awsx/awsx-common/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
)

var sessionName = util.RandomString(10)

func GetSessionWithAssumeRole(auth Auth) *session.Session {
	sess, err := awssession.GetSessionByCreds(auth.Region, auth.AccessKey, auth.SecretKey, "")

	if err != nil {
		fmt.Printf("failed to create aws session, %v\n", err)
		log.Fatal(err)
	}

	svc := sts.New(sess)

	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String(auth.CrossAccountRoleArn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int64(60 * 60 * 1),
	}

	if auth.ExternalId != "nil" {
		fmt.Println("Trying to fetch external id to assume new role")
		assumeRoleInput.ExternalId = aws.String(auth.ExternalId)
	}

	result, err := svc.AssumeRole(&assumeRoleInput)

	if err != nil {
		fmt.Printf("failed to assume role, %v\n", err)
		log.Fatal(err)
	}

	awsSession, err := awssession.GetSessionByCreds(auth.Region, *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken)

	if err != nil {
		fmt.Printf("failed to assume role, %v\n", err)
		log.Fatal(err)
	}

	return awsSession
}
