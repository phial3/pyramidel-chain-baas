// Copyright (c) 2022 s1ren
// hxx258456/pyramidel-chain-baas is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//
//	http://license.coscl.org.cn/MulanPSL2
//
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package localconfig

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	if err := os.Setenv("PYCBAAS_CFG_PATH", "e:\\github.com\\hxx258456\\pyramidel-chain-baas\\configs"); err != nil {
		t.Error(err)
	}
	loadConfig()
}
