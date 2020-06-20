package sshutils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sshShortDesc = "Provides ssh specific tooling"
	sshLongDesc  = `Provides added capability for ssh related stuff.
		For example it can generate cloudformation templates etc.`
	sshExample = `
	### Available commands for aws
	gotils ssh (connect)`
	ssh_private_key_path string
	ssh_public_key_path  string
	ssh_config_path      string
	ssh_username         string
	ssh_host             string
	ssh_port             int
)

func InitSsh() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "ssh [sub]",
		Short:            sshShortDesc,
		Long:             sshLongDesc,
		Example:          sshExample,
		TraverseChildren: true,
	}

	cmd.PersistentFlags().StringVarP(&ssh_private_key_path, "ssh-private-key-path", "", "", "ssh private key path")
	cmd.PersistentFlags().StringVarP(&ssh_public_key_path, "ssh-public-key-path", "", "", "ssh public key path")
	cmd.PersistentFlags().StringVarP(&ssh_config_path, "ssh-config-path", "", "", "ssh config path")
	cmd.PersistentFlags().StringVarP(&ssh_username, "ssh-username", "", "", "ssh username")
	viper.BindPFlag("ssh.private_key_path", cmd.PersistentFlags().Lookup("ssh-private-key-path"))
	viper.BindPFlag("ssh.public_key_path", cmd.PersistentFlags().Lookup("ssh-public-key-path"))
	viper.BindPFlag("ssh.config_path", cmd.PersistentFlags().Lookup("ssh-config-path"))
	viper.BindPFlag("ssh.username", cmd.PersistentFlags().Lookup("ssh-username"))
	cmd.AddCommand(NewSSHConnectCmd())
	return cmd
}
