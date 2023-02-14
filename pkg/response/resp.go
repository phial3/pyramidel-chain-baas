package response

type Response struct {
	Code int         `json:"code" ` // 响应码
	Msg  string      `json:"msg"`   // 消息
	Data interface{} `json:"data"`  // 数据
}
