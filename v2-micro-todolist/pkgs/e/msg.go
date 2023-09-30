package e

var MsgFlags = map[int]string{
	SUCCESS: "ok",
	Error:   "fail",

	InvalidParams: "请求参数错误",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}
