package sshutils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

// NewAwsConn - Setup a authenticated aws session for a specific assumed role
func NewAwsConn() *AwsConn {
	sess := awsSession.Must(awsSession.NewSession())
	creds := stscreds.NewCredentials(sess, "arn:aws:iam::302586665182:role/MasterAccountAdminAccessRole", func(p *stscreds.AssumeRoleProvider) {
		p.SerialNumber = aws.String(ssh_aws_token_serial_number)
		p.TokenProvider = stscreds.StdinTokenProvider
	})
	cfg := &aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: creds,
	}
	return &AwsConn{
		Session: sess,
		Config:  cfg,
	}
}

// GetUser - Get the information for the aws user which is configured for the current session
func (a *AwsConn) GetUser() {
	svc := iam.New(a.Session, a.Config)
	input := &iam.GetUserInput{}
	r, err := svc.GetUser(input)
	if err != nil {
		utils.CheckErr(fmt.Sprintf("Could not fetch user from aws - %s", err))
	}
	fmt.Println(r)
	return
}
