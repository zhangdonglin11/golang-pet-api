package sms

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"math/rand"
	"strings"
	"time"
)

func CreateClient(accessKeyID, accessKeySecret string) (_result *openapi.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(accessKeyID),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(accessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &openapi.Client{}
	_result, _err = openapi.NewClient(config)
	return _result, _err
}

func CreateApiInfo() (_result *openapi.Params) {
	params := &openapi.Params{
		// 接口名称
		Action: tea.String("SendSms"),
		// 接口版本
		Version: tea.String("2017-05-25"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
	_result = params
	return _result
}

// GenerateAnSMSCode 生成width长度的验证码
func GenerateAnSMSCode(width int) string {
	numerical := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numerical)
	random := rand.New(rand.NewSource(time.Now().Unix()))

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numerical[random.Intn(r)])
	}
	return sb.String()
}
