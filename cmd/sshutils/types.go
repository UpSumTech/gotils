package sshutils

import (
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
}
