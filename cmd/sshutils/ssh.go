package sshutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
	"golang.org/x/crypto/ssh"
)

var (
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stdio  bytes.Buffer
)

// NewSSHConnectCmd - is a function that generates a command to create a new connection and assign a pty
func NewSSHConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect SERVER PORT",
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
			ssh_host = args[0]
			if len(args) == 1 {
				ssh_port = 22
			}
			if len(args) == 2 {
				port, err := strconv.Atoi(args[1])
				if err != nil {
					return utils.RaiseCmdErr(cmd, "port needs to be an integer")
				}
				ssh_port = port
			}
			if len(args) > 2 {
				return utils.RaiseCmdErr(cmd, "Too many args")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			r, err := Terminal()
			if err != nil {
				utils.CheckErr(fmt.Sprintf("assigning pseudo terminal failed: %s", err))
			}
			Start(r)
		},
	}
	return cmd
}

// PublicKey - function to get the public key by reading the private key
func PublicKey() (ssh.AuthMethod, error) {
	buf, err := ioutil.ReadFile(ssh_private_key_path)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

// Config - returns a valid ssh config to dial the server with
func Config() (*ssh.ClientConfig, error) {
	pubkey, err := PublicKey()
	if err != nil {
		return nil, err
	}
	cfg := &ssh.ClientConfig{
		User: ssh_username,
		Auth: []ssh.AuthMethod{
			pubkey,
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	return cfg, nil
}

// Client - returns a valid ssh client
func Client() (*ssh.Client, error) {
	cfg, err := Config()
	if err != nil {
		return nil, err
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ssh_host, ssh_port), cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Session - returns a valid ssh session
func Session() (*ssh.Session, error) {
	conn, err := Client()
	if err != nil {
		return nil, err
	}
	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	return session, nil
}

// Terminal - create a new interactive terminal
func Terminal() (*RemoteShellConn, error) {
	session, err := Session()
	if err != nil {
		return nil, err
	}
	t := &TerminalConfig{
		Terminal: "xterm",
		Height:   80,
		Width:    40,
		Modes: ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		},
	}
	r := RemoteShellConn{
		Session:        session,
		TerminalConfig: t,
	}
	return &r, nil
}

// Start - starts an interactive terminal session and wait
func Start(r *RemoteShellConn) {
	fmt.Println("Starting remote shell")
	fmt.Println(r)
	if err := r.Session.RequestPty(r.TerminalConfig.Terminal, r.TerminalConfig.Height, r.TerminalConfig.Width, r.TerminalConfig.Modes); err != nil {
		r.Session.Close()
		utils.CheckErr(fmt.Sprintf("request for pseudo terminal failed: %s", err))
	}

	if err := r.Session.Shell(); err != nil {
		r.Session.Close()
		utils.CheckErr(fmt.Sprintf("request for opening a remote shell: %s", err))
	}

	if err := r.Session.Wait(); err != nil {
		r.Session.Close()
		utils.CheckErr(fmt.Sprintf("remote shell did not wait: %s", err))
	}
}
