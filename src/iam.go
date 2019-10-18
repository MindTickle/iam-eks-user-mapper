package main

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/kataras/golog"
)

func getAwsIamGroup(groupName string) *iam.GetGroupOutput {
	sess := session.Must(session.NewSession())
	config := &aws.Config{}

	if crossAccountRoleArn, exists := os.LookupEnv("AWS_IAM_ACCOUNT_ROLE_ARN"); exists {
		config.Credentials = stscreds.NewCredentials(sess, crossAccountRoleArn, func(arp *stscreds.AssumeRoleProvider) {
			arp.RoleSessionName = "devops-iam-eks-mapper-session"
			arp.Duration = 60 * time.Minute
			arp.ExpiryWindow = 30 * time.Second
		})
	}

	iamClient := iam.New(sess, config)
	group, err := iamClient.GetGroup(&iam.GetGroupInput{
		GroupName: aws.String(groupName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				golog.Error(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				golog.Error(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				golog.Error(aerr.Error())
			}
		}
	}
	return group
}
