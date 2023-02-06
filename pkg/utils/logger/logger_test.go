package logger

import (
	"github.com/hxx258456/pyramidel-chain-baas/internal/localconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"testing"
)

func TestClean(t *testing.T) {
	t.Cleanup(func() {
		pwd, err := os.Getwd()
		if err != nil {
			t.Error(err)
		}
		dir := path.Join(pwd, "testdata")
		if err = os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
		mode := os.FileMode(0644)

		if err = os.Mkdir(dir, mode); err != nil {
			t.Error(err)
		}
	})
}

func TestInfoLog(t *testing.T) {
	if err := os.Setenv("PYCBAAS_CFG_PATH", "E:\\github.com\\hxx258456\\pyramidel-chain-baas\\configs"); err != nil {
		t.Error(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	localconfig.Defaultconfig.Logger.Filename = path.Join(pwd, "testdata", "log_test.log")
	InitLogger(&localconfig.Defaultconfig.Logger)
	lg.Sugar().Info("this is a test")
}

func TestNamed(t *testing.T) {
	if err := os.Setenv("PYCBAAS_CFG_PATH", "E:\\github.com\\hxx258456\\pyramidel-chain-baas\\configs"); err != nil {
		t.Error(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	localconfig.Defaultconfig.Logger.Filename = path.Join(pwd, "testdata", "log_test.log")
	InitLogger(&localconfig.Defaultconfig.Logger)
	testLog := zap.L().Named("test")
	testLog.Sugar().Info("this is a test")
}

func TestWithFileds(t *testing.T) {
	if err := os.Setenv("PYCBAAS_CFG_PATH", "E:\\github.com\\hxx258456\\pyramidel-chain-baas\\configs"); err != nil {
		t.Error(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	localconfig.Defaultconfig.Logger.Filename = path.Join(pwd, "testdata", "log_test.log")
	InitLogger(&localconfig.Defaultconfig.Logger)
	testLog := zap.L().Named("test")
	fields := []zapcore.Field{
		zap.Int("status", 500),
		zap.String("method", "GET"),
	}
	testLog.Info("fileds:", fields...)

	dogLog := zap.L().Named("dog")

	dogLog.Info("fileds:", fields...)
}
