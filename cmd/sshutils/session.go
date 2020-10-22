package sshutils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
)

// Session - returns a valid ssh session
func NewSession(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*ssh.Session, error) {
	conn, err := Client(ssh_public_key_path, ssh_username, ssh_host, ssh_port)
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

func NewSshSessionSubprocess(stdin io.Reader, stdout io.Writer, stderr io.Writer, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command("ssh", args...)

	if stdin == nil {
		stdin = os.Stdin
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(cmdStdin, stdin) // Pipe to destination from source

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(stdout, cmdStdout) // Pipe to destination from source

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(stderr, cmdStderr) // Pipe to destination from source

	// ignore signal(sigint)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-sigs:
			case <-done:
				break
			}
		}
	}()
	defer close(done)

	return cmd, nil
}
