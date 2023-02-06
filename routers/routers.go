package routers

import "github.com/gin-gonic/gin"

type Option func(group *gin.RouterGroup)

var options []Option

// Include 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init(e *gin.Engine) {
	v0Router := e.Group("/v0")

	for _, opt := range options {
		opt(v0Router)
	}
}
