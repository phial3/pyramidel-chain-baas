// Copyright (c) 2022 s1ren
// hxx258456/pyramidel-chain-baas is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 			http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package logger

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"sync"
	"testing"
)

func TestInfoLog(t *testing.T) {
	CfgConsoleLogger(false, false)

	Info("can see me")
	Debug("cannot see me")

	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	filepath := path.Join(pwd, "testdata", "log_test.log")
	CfgConsoleAndFileLogger(true, filepath, true)

	Info("can see me")
	Debug("cannot see me")
	Warn("this is warn")
	Error("this is error: %s", errors.New("this is error"))
	Error("info %% is dead", errors.New("this is error"), 2)
	Error(errors.New("this is error"))
	Error(errors.New("this is error"), "more error")
	// i := 0
	// for i < 10000 {
	// 	Info("can see me")
	// 	Debug("cannot see me")
	// 	Warn("this is warn")
	// 	Error("this is error: %s", errors.New("this is error"))
	// 	Error("info %% is dead", errors.New("this is error"), 2)
	// 	Error(errors.New("this is error"))
	// 	Error(errors.New("this is error"), "more error")
	// 	i++
	// }

	if IsDebugMode() == false {
		t.Error("not in debug mode")
	}
}

func TestFatalLog(t *testing.T) {
	if os.Getenv("LOG_FATAL") == "1" {
		Fatal("this is fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalLog")
	cmd.Env = append(os.Environ(), "LOG_FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestPanicLog(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Panic("this panics")
}

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

func TestNewLogger(t *testing.T) {

	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	filepath := path.Join(pwd, "testdata", "log_test.log")
	CfgConsoleAndFileLogger(true, filepath, true)

	// 模块化日志测试
	logging := NewLogger("[测试]")
	fablogging := NewLogger("[fab测试]")
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {

		i := 0
		for i < 100 {
			logging.Info("can see me")
			logging.Debug("cannot see me")
			logging.Warn("this is warn")
			logging.Error("this is error: %s", errors.New("this is error"))
			logging.Error("info %% is dead", errors.New("this is error"), 2)
			logging.Error(errors.New("this is error"))
			logging.Error(errors.New("this is error"), "more error")
			i++
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		i := 0
		for i < 100 {
			fablogging.Info("can see me")
			fablogging.Debug("cannot see me")
			fablogging.Warn("this is warn")
			fablogging.Error("this is error: %s", errors.New("this is error"))
			fablogging.Error("info %% is dead", errors.New("this is error"), 2)
			fablogging.Error(errors.New("this is error"))
			fablogging.Error(errors.New("this is error"), "more error")
			i++
		}
		wg.Done()
	}()
	wg.Wait()
}
