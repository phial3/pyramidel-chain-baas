package member

type DownloadKS struct {
	Uscc string `json:"uscc" binding:"required"`
	Name string `json:"name" binding:"required"`
} // 下载用户私钥接口

type UpdateFrozen struct {
	Uscc     string `json:"uscc" binding:"required"`
	Name     string `json:"name" binding:"required"`
	IsFrozen *bool  `json:"isfrozen" binding:"required"`
} // 更新用户冻结状态

type RegenerateTokenReq struct {
	Uscc string `json:"uscc" binding:"required"`
	Name string `json:"name" binding:"required"`
} // 重新生成证书

type revoke struct {
	Uscc string `json:"uscc" binding:"required"`
	Name string `json:"name" binding:"required"`
} // 注销用户
