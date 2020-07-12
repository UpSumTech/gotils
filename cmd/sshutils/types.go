package sshutils

import (
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/crypto/ssh"
)

// TerminalConfig - configuration for the ssh pty
type TerminalConfig struct {
	Terminal string
	Height   int
	Width    int
	Modes    ssh.TerminalModes
}

// RemoteShellConn - configuration for remote shell connection
type RemoteShellConn struct {
	Session        *ssh.Session
	TerminalConfig *TerminalConfig
	InteractivePty bool
}

type AwsConn struct {
	Session *awsSession.Session
	Config  *aws.Config
}
