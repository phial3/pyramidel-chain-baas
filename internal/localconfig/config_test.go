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
