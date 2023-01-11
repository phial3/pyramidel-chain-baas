package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/hxx258456/pyramidel-chain-baas/pkg/constants"
	"github.com/spf13/viper"
)

func InitViper(v *viper.Viper, configName string) error {
	if v == nil {
		return errors.New("nil pointer")
	}
	var altPath = os.Getenv("PYCBAAS_CFG_PATH")
	if altPath != "" {
		// If the user has overridden the path with an envvar, its the only path
		// we will consider
		if !dirExists(altPath) {
			return fmt.Errorf("PYCBAAS_CFG_PATH %s does not exist", altPath)
		}

		v.AddConfigPath(altPath)
	} else {
		// If we get here, we should use the default paths in priority order:
		//

		// And finally, the official path
		if dirExists(constants.OfficialPath) {
			v.AddConfigPath(constants.OfficialPath)
		}
	}

	v.SetConfigName(configName)

	return nil
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
