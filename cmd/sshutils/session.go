package sshutils

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

// Session - returns a valid ssh session
func NewSession(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*ssh.Session, error) {
	conn, err := Client()
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	if stdin == nil {
		stdin = os.Stdin
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	sessionStdin, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(sessionStdin, stdin) // Pipe to destination from source

	sessionStdout, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(stdout, sessionStdout) // Pipe to destination from source

	sessionStderr, err := session.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(stderr, sessionStderr) // Pipe to destination from source

	return session, nil
}
