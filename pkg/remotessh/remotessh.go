package remotessh

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"github.com/melbahja/goph"
	probing "github.com/prometheus-community/pro-bing"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

var sshLogger = logger.Lg.Named("remote/ssh")

func ConnectToHost(username, pass, ip string, port uint) (*goph.Client, error) {
	client, err := goph.NewConn(&goph.Config{
		Auth:     goph.Password(pass),
		User:     username,
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
	// here you should check new connections public keys
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

func CheckCommand(command string, client *goph.Client) (bool, error) {
	cmd := "which " + command
	_, err := client.Run(cmd)
	if err != nil {
		return false, err
	}
	return true, nil
}

func InitHost(client *goph.Client) (out []byte, ok bool, err error) {
	if err := client.Upload("/root/pyramidel-chain-baas/scripts/init.sh", "/root/init.sh"); err != nil {
		return nil, false, err
	}

	// 赋予init.sh执行权限,并执行
	_, err = client.Run("chmod +x ~/init.sh")
	if err != nil {
		return nil, false, err
	}
	initOut, err := client.Run("cd ~ && ./init.sh")
	if err != nil {
		return nil, false, err
	}
	sshLogger.Debug(string(initOut))

	return initOut, true, nil
}

func Ping(ip string) int64 {

	pinger, err := probing.NewPinger(ip)
	if err != nil {
		sshLogger.Error("Ping ip %s failed", zap.Error(err))
		return 0
	}
	//pinger.SetPrivileged(false)
	pinger.Count = 3
	pinger.Timeout = time.Second * 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		sshLogger.Error("Ping ip %s failed", zap.Error(err))
		return 0
	}
	stats := pinger.Statistics()
	return stats.AvgRtt.Microseconds()
}

func EnrollBootstrap(client *goph.Client, url, caname string) error {
	cmd, err := client.Command("fabric-ca-client", "enroll", "-d", "-u", url, "--caname", caname, "--tls.certfiles")
	if err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
