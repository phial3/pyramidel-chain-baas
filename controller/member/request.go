package member

type DownloadKS struct {
	Uscc string `json:"uscc" binding:"required"`
	Name string `json:"name" binding:"required"`
} // 下载用户私钥接口
