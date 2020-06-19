package awsutils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	awsShortDesc = "Provides aws specific tooling"
	awsLongDesc  = `Provides added capability for aws related stuff.
		For example it can interact with aws resources through the api for convenient tasks`
	awsExample = `
	### Available commands for aws
	gotils aws (ec2|s3)`
	aws_access_key_id     string
	aws_secret_access_key string
	aws_account_id        string
	aws_region            string
)

func InitAws() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aws [sub]",
		Short:   awsShortDesc,
		Long:    awsLongDesc,
		Example: awsExample,
		PreRun: func(cmd *cobra.Command, args []string) {
			aws_access_key_id = viper.GetString("aws.access_key_id")
			aws_secret_access_key = viper.GetString("aws.secret_access_key")
			aws_account_id = viper.GetString("aws.account_id")
			aws_region = viper.GetString("aws.region")
		},
		TraverseChildren: true,
	}

	cmd.AddCommand(NewAwsEC2Cmds())
	cmd.Flags().StringVarP(&aws_access_key_id, "aws_access_key_id", "", "", "Aws access key id")
	cmd.Flags().StringVarP(&aws_secret_access_key, "aws_secret_access_key", "", "", "Aws secret access key")
	cmd.Flags().StringVarP(&aws_region, "aws_region", "", "", "Aws region")
	cmd.Flags().StringVarP(&aws_account_id, "aws_account_id", "", "", "Aws account id")
	viper.BindPFlag("aws.access_key_id", cmd.Flags().Lookup("aws_access_key_id"))
	viper.BindPFlag("aws.secret_access_key", cmd.Flags().Lookup("aws_secret_access_key"))
	viper.BindPFlag("aws.account_id", cmd.Flags().Lookup("aws_account_id"))
	viper.BindPFlag("aws.region", cmd.Flags().Lookup("aws_region"))
	return cmd
}
