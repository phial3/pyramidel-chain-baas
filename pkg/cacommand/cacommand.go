// Copyright (c) 2022 s1ren
// hxx258456/pyramidel-chain-baas is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 			http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package cacommand

import (
	"bytes"
	"errors"
	"go.uber.org/zap"
	"os/exec"
)

var logging = zap.L().Named("cacommand")

func EnrollBootstrap(url, caname, certfiles string) error {
	cmd := exec.Command("fabric-ca-client", "enroll", "-u", url, "--caname", caname, "--tls.certfiles", certfiles)
	logging.Info(cmd.String())
	var stdOut, stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut
	if err := cmd.Run(); err != nil {
		return err
	}
	if stdErr.Bytes() != nil {
		return errors.New(stdErr.String())
	}
	logging.Info(stdOut.String())
	return nil
}
