package awsutils

import (
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	dest string
)

func NewAwsEC2Cmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ec2 [sub]",
		Short:   "Provides ec2 utilities for aws",
		Long:    `Lets you filter and fetch ec2 resources for aws.
			Lets you perform different queries on ec2 instances and lets you filter them and fetch different ec2 resources.`,
		Example: `
			### Available commands for aws ec2 subcommands
			# gotils aws ec2 TEMPLATE_KIND
			gotils aws ec2 find `,
		TraverseChildren: true,
	}
	cmd.AddCommand(ec2FindByTagSubCmd())
	return cmd
}

func ec2FindByTagSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "findByTag TAG_NAME",
		Short:   "Find aws ec2 instances by tags",
		Long:    `Lets you filter and fetch ec2 resources for aws by tags.
			You can either provide tag name and/or tag value as well.`,
		Example: `
			### Available commands for aws ec2 subcommands
			# gotils aws ec2 TEMPLATE_KIND
			gotils aws ec2 find `,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return utils.RaiseCmdErr(cmd, "tag name and tag value to filter ec2 instances not provided")
			}
			if len(args) > 2 {
				return utils.RaiseCmdErr(cmd, "Too many args")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ec2FindByTag(args[0], args[1])
		},
	}

	return cmd
}

func ec2FindByTag(tag string, val string) {
	ec2svc := ec2.New(session.New())
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", tag)),
				Values: []*string{aws.String(val)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	for idx, res := range resp.Reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID: ", *inst.InstanceId)
		}
	}
}
