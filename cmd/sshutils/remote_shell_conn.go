package sshutils

import (
	"fmt"
	"os"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
	"golang.org/x/crypto/ssh"
)

// NewRemoteShellConn - create a new remote shell connection object
func NewRemoteShellConn(interactive bool, sshWithSsm bool) (*RemoteShellConn, error) {
	if sshWithSsm {
		awsConn := NewAwsConn()
		awsConn.GetUser()
	}
	session, err := NewSession(os.Stdin, os.Stdout, os.Stderr)
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
		InteractivePty: interactive,
	}

	return &r, nil
}

// StartInteractiveShell - starts an interactive terminal session and wait
func (r *RemoteShellConn) StartInteractiveShell() {
	if !r.InteractivePty {
		return
	}

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