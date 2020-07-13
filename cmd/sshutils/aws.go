package sshutils

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// NewAwsConn - Setup a authenticated aws session for a specific assumed role
func NewAwsConn() *AwsConn {
	// You need this for entropy of the random string generation
	rand.Seed(time.Now().UnixNano())

	userName, ssmSessionTagVal := GetAwsUserInfo()
	fmt.Println(userName)

	sess := awsSession.Must(awsSession.NewSession())

	creds := stscreds.NewCredentials(sess, "arn:aws:iam::302586665182:role/MasterAccountBastionSessionManagerAccessRole", func(p *stscreds.AssumeRoleProvider) {
		p.SerialNumber = aws.String(ssh_aws_token_serial_number)
		p.TokenProvider = stscreds.StdinTokenProvider
		p.RoleSessionName = ssh_host + "-" + ssmSessionTagVal + "-" + getRandString(20)
		p.ExpiryWindow = time.Duration(10) * time.Minute
		p.Tags = []*sts.Tag{
			&sts.Tag{
				Key:   aws.String("SSMSessionRunAs"),
				Value: aws.String(ssmSessionTagVal),
			},
			&sts.Tag{
				Key:   aws.String(ssh_host),
				Value: aws.String(ssh_host),
			},
		}
	})

	cfg := &aws.Config{
		Region:                        aws.String("us-west-2"),
		Credentials:                   creds,
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	return &AwsConn{
		Session: sess,
		Config:  cfg,
	}
}

// GetBastionInstance - fetches the bastion instance to log into based on tags
func (a *AwsConn) GetBastionInstance() string {
	var instanceId string
	ec2svc := ec2.New(a.Session, a.Config)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", "Service")),
				Values: []*string{aws.String("bastion")},
			},
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", "Environment")),
				Values: []*string{aws.String("dev")},
			},
		},
	}

	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			instanceId = *inst.InstanceId
		}
	}

	return instanceId
}

// GetAwsUserInfo - Get the information for the aws user which is configured for the current session
func GetAwsUserInfo() (string, string) {
	sess := awsSession.Must(awsSession.NewSession())
	cfg := &aws.Config{
		Region: aws.String("us-west-2"),
	}
	svc := iam.New(sess, cfg)
	input := &iam.GetUserInput{}
	r, err := svc.GetUser(input)
	if err != nil {
		utils.CheckErr(fmt.Sprintf("Could not fetch user from aws - %s", err))
	}
	userName := *r.User.UserName
	ssmSessionTagVal := findSsmSessionTagVal(r.User.Tags)
	return userName, ssmSessionTagVal
}

func findSsmSessionTagVal(tags []*iam.Tag) string {
	var res string
	for _, tag := range tags {
		if *tag.Key == "SSMSessionRunAs" {
			res = *tag.Value
		}
	}
	return res
}

func getRandString(n int) string {
	buf := make([]rune, n)
	for i := 0; i < n; i++ {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}
