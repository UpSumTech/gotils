package sshutils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

// NewSSHConnectCmd - is a function that generates a command to create a new connection and assign a pty
func NewSSHConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect SERVER",
		Short: "Connects to a remote server",
		Long: `Lets you connect to a remote server.
			And assigns a pseudo terminal`,
		Example: `
			### Available commands for ssh subcommands
			gotils ssh connect SERVER`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(ssh_username) == 0 {
				ssh_username = viper.GetString("ssh.username")
			}
			if len(ssh_private_key_path) == 0 {
				ssh_private_key_path = viper.GetString("ssh.private_key_path")
			}
			if len(ssh_public_key_path) == 0 {
				ssh_public_key_path = viper.GetString("ssh.public_key_path")
			}
			if len(ssh_config_path) == 0 {
				ssh_config_path = viper.GetString("ssh.config_path")
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return utils.RaiseCmdErr(cmd, "server needs to be provided")
			}
			if len(args) > 1 {
				return utils.RaiseCmdErr(cmd, "Too many args")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ssh_host = args[0]
			sshConnect()
		},
	}
	return cmd
}

func sshConnect() {
	str := fmt.Sprintf("Using ssh username %s\n", ssh_username)
	str += fmt.Sprintf("Using private key %s\n", ssh_private_key_path)
	str += fmt.Sprintf("Using public key %s\n", ssh_public_key_path)
	str += fmt.Sprintf("Using ssh config %s\n", ssh_config_path)
	str += fmt.Sprintf("Connecting %s to %s\n", ssh_username, ssh_host)
	fmt.Println(str)
}
