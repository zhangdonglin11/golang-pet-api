// 状态码
package result

// codes定义
type Codes struct {
	SUCCESS        uint
	FAILED         uint
	NOAUTH         uint
	AUTHFORMATEROR uint

	Message map[uint]string
}

// ApiCode 状态码
var ApiCode = &Codes{
	SUCCESS:        200,
	FAILED:         501,
	NOAUTH:         403,
	AUTHFORMATEROR: 405,
}

// 状态信息
func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.SUCCESS:        "请求成功",
		ApiCode.FAILED:         "请求失败",
		ApiCode.NOAUTH:         "请求头中的auth为空",
		ApiCode.AUTHFORMATEROR: "请求头中的auth格式错误",
	}
}

// 供外部调用
func (c *Codes) GetMessage(code uint) string {
	message, ok := c.Message[code]
	if !ok {
		return ""
	}
	return message
}
