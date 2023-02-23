// Copyright (c) 2022 s1ren
// hxx258456/pyramidel-chain-baas is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 			http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package localconfig

import (
	"github.com/fsnotify/fsnotify"
	coreconfig "github.com/hxx258456/pyramidel-chain-baas/pkg/utils/config"
	"github.com/spf13/viper"
	"log"
)

var Defaultconfig = TopLevel{
	Logger: Logger{
		MaxAge:     15,
		MaxBackups: 3,
		MaxSize:    500,
		Level:      "info",
		Compress:   true,
		LocalTime:  true,
	},
	Serve: Serve{
		Mode: "debug",
		Port: ":8080",
		IpWhiteList: []string{
			"127.0.0.1",
			"localhost",
			"192.168.0.1",
			"172.17.0.1",
		},
	},
	Mysql: Mysql{
		Host:      "localhost",
		Port:      3306,
		User:      "root",
		Password:  "123456",
		Parsetime: true,
		Loc:       true,
		Charset:   "utf8_unicode_ci",
	},
	AMQP: AMQP{
		Host:        "47.92.54.239",
		Port:        5672,
		User:        "txhy",
		Password:    "txhy2022.com",
		Vhost:       "//auth",
		Queue:       "baasOrgAdd",
		ContentType: "application/json",
	},
}
var config *viper.Viper

type (
	TopLevel struct {
		Logger Logger `json:"logger" yaml:"logger"`
		Serve  Serve  `json:"serve" yaml:"serve"`
		Mysql  Mysql  `json:"mysql" yaml:"mysql"`
		AMQP   AMQP   `json:"amqp" yaml:"amqp"`
	}
)

func init() {
	//if err := os.Setenv("PYCBAAS_CFG_PATH", "E:\\github.com\\hxx258456\\pyramidel-chain-baas\\configs"); err != nil {
	//	panic(err)
	//}
	//加载配置
	loadConfig()
	watchConfig()
}

// 监听配置改变
func watchConfig() {
	go func() {
		config.WatchConfig()
		config.OnConfigChange(func(e fsnotify.Event) {

			//改变重新加载
			loadConfig()
		})
	}()
}

// 加载配置
func loadConfig() {
	config = viper.New()

	if err := coreconfig.InitViper(config, "config"); err != nil {
		panic(err)
	}

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	err := config.Unmarshal(&Defaultconfig)
	if err != nil {
		panic(err)
	}

	Defaultconfig.check()
	log.Printf("%+v", Defaultconfig)
	log.Println()
}

// checkConfig 检查配置格式
func (t *TopLevel) check() {
	t.Serve.check()
	t.Logger.check()
}
