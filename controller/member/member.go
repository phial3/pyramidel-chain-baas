package member

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/ccgo/sm2"
	"github.com/hxx258456/ccgo/x509"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/gmtoken"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/remotessh"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

var member = newMemberController()

type memberController struct {
}

func newMemberController() *memberController {
	return &memberController{}
}

func (c *memberController) New(ctx *gin.Context) {
	param := &model.Member{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		response.Error(ctx, err)
		log.Println("31:::::::::::::::::::::Error: ", err)
		return
	}

	// 申请证书
	var result []model.Organization
	org := model.Organization{}
	if err := org.FindByUscc(param.Uscc, &result); err != nil {
		response.Error(ctx, err)
		return
	}

	if len(result) <= 0 {
		response.Error(ctx, errors.New("organization not found"))
		return
	}

	param.OrganizationId = result[0].ID

	switch param.UserType {
	case "client":
		sshcli, err := remotessh.ConnectToHost(result[0].Host.Username, result[0].Host.Pw, result[0].Host.UseIp, result[0].Host.SSHPort)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		defer sshcli.Close()
		if err := remotessh.RegisterUser(sshcli, param.Uscc, param.Name, param.PassWord, "client", strconv.Itoa(int(result[0].CaServerPort))); err != nil {
			response.Error(ctx, err)
			return
		}
	case "admin":
		sshcli, err := remotessh.ConnectToHost(result[0].Host.Username, result[0].Host.Pw, result[0].Host.UseIp, result[0].Host.SSHPort)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		defer sshcli.Close()
		if err := remotessh.RegisterUser(sshcli, param.Uscc, param.Name, param.PassWord, "admin", strconv.Itoa(int(result[0].CaServerPort))); err != nil {
			response.Error(ctx, err)
			return
		}
	default:
		response.Error(ctx, errors.New("invalid user type"))
		return
	}

	switch param.StoreType {
	case 2:
		// 提供密钥证书下载链接
		if err := param.Create(); err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, nil, "enroll user success")
		return
	case 1:
		// 根据申请到的证书生成token返回给用户
		// 提供私钥下载链接
		payloads := gmtoken.CreateStdPayloads(param.Name, "test", "anyone", "N0001", 10*365*24*60)
		mspdir := fmt.Sprintf("/root/txhyjuicefs/organizations/%s/users/%s@%s.pcb.com/msp/", param.Uscc, param.Name, param.Uscc)
		keyDir := filepath.Join(mspdir, "keystore")
		// there's a single file in this dir containing the private key
		files, err := ioutil.ReadDir(keyDir)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		if len(files) != 1 {
			response.Error(ctx, fmt.Errorf("keystore folder should have contain one file"))
			return
		}
		keyPath := filepath.Join(keyDir, files[0].Name())

		privKey, err := x509.ReadPrivateKeyFromPemFile(keyPath, nil)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		token, err := gmtoken.BuildTokenWithGM(payloads, time.Time{}, privKey.(*sm2.PrivateKey))
		if err != nil {
			response.Error(ctx, err)
			return
		}
		param.Token = token
		if err := param.Create(); err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, gin.H{"token": token}, "enroll user success")
		return
	default:
		response.Error(ctx, errors.New("invalid store type"))
		return
	}
}

func (c *memberController) DownloadKeyStore(ctx *gin.Context) {
	param := &DownloadKS{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		response.Error(ctx, err)
		log.Println("31:::::::::::::::::::::Error: ", err)
		return
	}

	mspdir := fmt.Sprintf("/root/txhyjuicefs/organizations/%s/users/%s@%s.pcb.com/msp/", param.Uscc, param.Name, param.Uscc)
	keyDir := filepath.Join(mspdir, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if len(files) != 1 {
		response.Error(ctx, fmt.Errorf("keystore folder should have contain one file"))
		return
	}

	keyPath := filepath.Join(keyDir, files[0].Name())
	ctx.File(keyPath)
	return
}

func (c *memberController) DownloadCert(ctx *gin.Context) {
	param := &DownloadKS{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		response.Error(ctx, err)
		return
	}

	certPwd := fmt.Sprintf("/root/txhyjuicefs/organizations/%s/users/%s@%s.pcb.com/msp/signcerts/cert.pem", param.Uscc, param.Name, param.Uscc)
	ctx.File(certPwd)
	return
}
