package sshutils

import (
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

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
