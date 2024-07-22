// 验证码 服务层
package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store base64Captcha.Store = RedisStore{}

// 生成图片验证码
func CaptMake(len int) (id, b64s string, err error) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          len,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, store)
	lid, lb64s, _, lerr := captcha.Generate()
	return lid, lb64s, lerr
}
