package remotessh

import (
	"github.com/melbahja/goph"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

var sshLogger = zap.L().Named("remote/ssh")

func ConnectToHost(pass, ip string, port uint) (*goph.Client, error) {
	client, err := goph.NewConn(&goph.Config{
		Auth:     goph.Password(pass),
		User:     "root",
		Addr:     ip,
		Port:     port,
		Timeout:  10 * time.Second,
		Callback: VerifyHost,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {

	//
	// If you want to connect to new hosts.
	// here your should check new connections public keys
	// if the key not trusted you shuld return an error
	//

	// hostFound: is host in known hosts file.
	// err: error if key not in known hosts file OR host in known hosts file but key changed!
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {

		return err
	}

	// handshake because public key already exists.
	if hostFound && err == nil {

		return nil
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}
