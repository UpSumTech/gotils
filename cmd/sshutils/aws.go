package sshutils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/mitchellh/go-homedir"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func NewRootAwsSession(profile string, region string) (*awsSession.Session, error) {
	sess := awsSession.Must(awsSession.NewSessionWithOptions(awsSession.Options{
		SharedConfigState: awsSession.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
	}))
	_, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func NewAssumeRoleAwsSession(tags []*sts.Tag) (*awsSession.Session, error) {
	// You need this for entropy of the random string generation
	rand.Seed(time.Now().UnixNano())

	sess, err := NewRootAwsSession("work_root", "us-west-2")
	if err != nil {
		return nil, err
	}

	// Generate sts creds with ssm session tags
	creds := stscreds.NewCredentials(sess, "arn:aws:iam::302586665182:role/MasterAccountBastionSessionManagerAccessRole", func(p *stscreds.AssumeRoleProvider) {
		p.SerialNumber = aws.String(ssh_aws_token_serial_number)
		p.TokenProvider = stscreds.StdinTokenProvider
		p.RoleSessionName = ssh_host + "-" + getRandString(20)
		p.ExpiryWindow = time.Duration(10) * time.Minute
		p.Tags = tags
	})

	// Generate aws config for service calls with the temp sts creds
	newSession := awsSession.Must(awsSession.NewSessionWithOptions(awsSession.Options{
		SharedConfigState: awsSession.SharedConfigEnable,
		Config: aws.Config{
			Region:                        aws.String("us-west-2"),
			Credentials:                   creds,
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}))

	UpdateAwsSharedCredentials(newSession)

	return newSession, nil
}

func UpdateAwsSharedCredentials(sess *awsSession.Session) {
	home, err := homedir.Dir()
	if err != nil {
		utils.CheckErr(fmt.Sprintf("could not find home dir path", err))
	}

	file, err := os.OpenFile(filepath.Join(home, ".aws/credentials"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		utils.CheckErr(fmt.Sprintf("could not open aws credentials file to update profile - %s", err))
	}
	defer file.Close()

	val, err := sess.Config.Credentials.Get()
	if err != nil {
		utils.CheckErr(fmt.Sprintf("could not retrieve credentials of assumed role session - %s", err))
	}

	text := fmt.Sprintf(`
[%s]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token = %s
`, "gotils_ssh", val.AccessKeyID, val.SecretAccessKey, val.SessionToken)

	_, err = file.WriteString(text)
	if err != nil {
		utils.CheckErr(fmt.Sprintf("could not update aws credentials file to update profile - %s", err))
	}
}

// NewAwsConn - Setup a authenticated aws session for a specific assumed role
func NewAwsConn() (*AwsConn, error) {
	tags, err := GetSessionTags()
	if err != nil {
		return nil, err
	}
	sess, err := NewAssumeRoleAwsSession(tags)
	if err != nil {
		return nil, err
	}
	return &AwsConn{
		Session: sess,
	}, nil
}

// GetSshTargetInstance - fetches the ec2 instance to log into based on tags
func (a *AwsConn) GetSshTargetInstance() string {
	var instanceId string
	ec2svc := ec2.New(a.Session)

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

func GetSessionTags() ([]*sts.Tag, error) {
	_, sessionTagVal, err := GetAwsUserInfo()
	if err != nil {
		return nil, err
	}

	tags := []*sts.Tag{
		&sts.Tag{
			Key:   aws.String("SSMSessionRunAs"),
			Value: aws.String(sessionTagVal),
		},
		&sts.Tag{
			Key:   aws.String(ssh_host),
			Value: aws.String(ssh_host),
		},
	}
	return tags, nil
}

// GetAwsUserInfo - Get the information for the aws user which is configured for the current session
func GetAwsUserInfo() (string, string, error) {
	sess, err := NewRootAwsSession("work_root", "us-west-2")
	if err != nil {
		return "", "", err
	}
	svc := iam.New(sess)
	input := &iam.GetUserInput{}
	r, err := svc.GetUser(input)
	if err != nil {
		return "", "", err
	}
	userName := *r.User.UserName
	ssmSessionTagVal := findSsmSessionTagVal(r.User.Tags)
	return userName, ssmSessionTagVal, nil
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
