// Copyright (c) 2022 s1ren
// hxx258456/pyramidel-chain-baas is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 			http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package localconfig

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
		// Filename 日志文件存放路径,默认为<processname>-lumberjack.log in os.TempDir()
		Filename string `json:"filename" yaml:"filename"`

		// MaxSize 单个日志文件最大大小，以MB为单位
		MaxSize int `json:"maxsize" yaml:"maxsize"`

		// MaxAge 保留日志文件最大天数以自然天为单位
		MaxAge int `json:"maxage" yaml:"maxage"`

		// MaxBackups 最大备份文件数量
		MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

		// LocalTime 是否使用计算机本地时间,默认使用UTC
		LocalTime bool `json:"localtime" yaml:"localtime"`

		// Compress 是否压缩文件
		// TODO: 添加压缩算法选择gzip g4
		Compress bool `json:"compress" yaml:"compress"`

		// SkipPaths logger中间件忽略的请求路由,无效路由会被抛弃
		SkipPaths []string `json:"skippaths,omitempty" yaml:"skippaths,omitempty"`

		// Context trace跟踪
		Context Fn `json:"-" yaml:"-"`

		// Level 日志级别info,INFO,error,ERROR
		Level string `json:"level" yaml:"level" default:"info"`
	}
)

type Fn func(c *gin.Context) []zapcore.Field

func (l *Logger) check() {}
