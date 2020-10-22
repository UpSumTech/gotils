package sshutils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/viper"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
	"golang.org/x/crypto/ssh"
)

// NewSsmShellConn - create a new remote shell connection object
func NewSsmShellConn(awsConn *AwsConn, interactive bool) (*exec.Cmd, error) {
	ssmPluginPath, err := exec.LookPath(viper.GetString("ssh.ssm_plugin_name"))
	if err != nil {
		return nil, err
	}
	fmt.Println("Using session manager plugin at : ", ssmPluginPath)

	instanceId := awsConn.GetSshTargetInstance()
	input := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartSSHSession"),
		Parameters: map[string][]*string{
			"portNumber": []*string{aws.String(strconv.Itoa(ssh_port))},
		},
		Target: aws.String(instanceId),
	}

	inputJson, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	ssmsvc := ssm.New(awsConn.Session)
	ssmSession, err := ssmsvc.StartSession(input)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := TerminateSsmSession(awsConn, *ssmSession.SessionId)
		if err != nil {
			utils.CheckErr(fmt.Sprintf("Failed to terminate ssm session: %s", err))
		}
	}()

	ssmSessionJson, err := json.Marshal(ssmSession)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(inputJson))
	fmt.Println(string(ssmSessionJson))
	endpoint := ssmsvc.Client.Endpoint
	proxyCmd := fmt.Sprintf("ProxyCommand=%s '%s' %s %s %s '%s' %s", "session-manager-plugin", string(ssmSessionJson), *awsConn.Session.Config.Region, "StartSession", "gotils_ssh", string(inputJson), endpoint)
	sshArgs := []string{"-i", ssh_private_key_path, "-tt", "-o", proxyCmd, instanceId}
	cmd, err := NewSshSessionSubprocess(nil, nil, nil, sshArgs...)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func TerminateSsmSession(awsConn *AwsConn, sessionId string) error {
	ssmsvc := ssm.New(awsConn.Session)
	_, err := ssmsvc.TerminateSession(&ssm.TerminateSessionInput{
		SessionId: &sessionId,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewSshShellConn(interactive bool) (*RemoteShellConn, error) {
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
