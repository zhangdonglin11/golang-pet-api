package v1

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-pet-api/global"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/utils"
	"golang-pet-api/utils/captcha"
	"golang-pet-api/utils/imageUtils"
	"golang-pet-api/utils/jwt"
	"golang-pet-api/utils/sms"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

type Auth struct {
}

// Captcha godoc
// @Summary 获取图片验证码
// @Description 获取一个基于数字和字母的图片验证码，并返回验证码的ID和Base64编码的图片数据。
// @Tags 认证模块
// @Produce json
// @Success 200 {string} json{Code,Msg,Data} "成功"
// @Router /api/v1/captcha [get]
func (a Auth) Captcha(c *gin.Context) {
	//生成验证码
	id, base64Image, err := captcha.CaptMake(5)
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
		return
	}
	data := map[string]interface{}{
		"captchaId": id,
		"captcha":   base64Image,
	}
	utils.RespOk(c.Writer, data, "请求成功")
}

// SmsCode godoc
// @Summary 获取短信验证码
// @Description 通过电话号码和类型获取短信验证码。
// @Tags 认证模块
// @Produce json
// @Param mobile formData string true "电话号码" format(phone) default(13169197369)
// @Param type formData string true "类型" Enums(0, 1) default(1)
// @Success 200 {string} json{Code,Msg,Data} "成功"
// @Router /api/v1/sms [post]
func (a Auth) SmsCode(c *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	//获取验证码节流
	ttl := global.RedisDb.TTL(c, sendSmsForm.Mobile)
	// 验证码有效300秒
	if ttl.Val() > time.Second*240 {
		message := fmt.Sprintf("%s秒后可获取验证码", ttl.Val()-time.Second*240)
		utils.RespFail(c.Writer, message)
		return
	}

	//生成验证码
	code := sms.GenerateAnSMSCode(5)
	//创建短信服务
	client, _err := sms.CreateClient(global.Config.Sms.AccessKeyID, global.Config.Sms.AccessKeySecret)
	if _err != nil {
		utils.RespFail(c.Writer, "短信服务："+_err.Error())
		return
	}
	params := sms.CreateApiInfo()
	// query params
	queries := map[string]interface{}{}
	queries["PhoneNumbers"] = tea.String(sendSmsForm.Mobile)
	queries["SignName"] = tea.String("宠物萌小程序")
	queries["TemplateCode"] = tea.String("SMS_464376212")
	//fmt.Sprintf("{\"code\":\"%s\"}", code)
	queries["TemplateParam"] = tea.String("{\"code\":\"" + code + "\"}")
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	if _, _err = client.CallApi(params, request, runtime); _err != nil {
		utils.RespFail(c.Writer, "短信服务："+_err.Error())
		return
	}
	if err := global.RedisDb.Set(c, sendSmsForm.Mobile, code, time.Second*300); err.Err() != nil {
		utils.RespFail(c.Writer, fmt.Sprintln("redis短信验证码:", err.Err()))
		return
	}
	utils.RespOk(c.Writer, sendSmsForm.Mobile, "获取成功")
	return
}

// Login godoc
// @Summary 登录 用户/密码
// @Description 用户名和密码登录
// @Tags 认证模块
// @Produce json
// @Param username formData string true "用户名" minlength(1) maxlength(50) default(z1)
// @Param password formData string true "用户密码" minlength(6) maxlength(16) default(123456)
// @Success 200 {string} json{Code,Msg,Data} "成功"
// @Router /api/v1/login [post]
func (a Auth) Login(c *gin.Context) {
	var from forms.LoginForm
	if err := c.ShouldBind(&from); err != nil {
		HandleValidatorError(c, err)
		return
	}
	var user model.User
	if result := global.Db.Where("username=? ", from.UserName).First(&user); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "用户名不存在")
		return
	}
	if !utils.VerifyPassword(from.Password, user.Salt, user.Password) {
		utils.RespFail(c.Writer, "用户名或密码错误")
		return
	}
	//生成token 设置token保存的信息
	data := map[string]interface{}{
		"userId": user.ID,
		"role":   user.Role,
	}
	token, err := jwt.GetJwtToken(jwt.TokenKey, time.Now().Unix(), jwt.TokenExpirationTime, data)
	if err != nil {
		utils.RespFail(c.Writer, "token生成失败！")
		return
	}
	utils.RespOk(c.Writer, token, "登录成功！")
	return
}

// Register godoc
// @Summary 注册 用户名/密码
// @Description 用户名和密码注册。
// @Tags 认证模块
// @Produce json
// @Param username formData string true "用户名" minlength(1) maxlength(50)
// @Param password formData string true "用户密码" minlength(6) maxlength(16)
// @Param captcha formData string true "验证码" length(5)
// @Param captchaId formData string true "图片验证码key"
// @Success 200 {string} json{Code,Msg,Data} "成功"
// @Router /api/v1/register [post]
func (a Auth) Register(c *gin.Context) {
	var form forms.RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}
	fmt.Println(form)

	// 验证图片验证码
	booled := captcha.RedisStore{}.Verify(form.CaptchaId, form.Captcha, false)
	if !booled {
		utils.RespFail(c.Writer, "验证码错误")
		return
	}

	// 查询用户是否已经注册
	var user model.User
	if result := global.Db.Where("user_id=?", form.UserName).First(&user); result.RowsAffected != 0 {
		utils.RespFail(c.Writer, "用户名已被注册")
		return
	}

	//初始化用户表
	rand.New(rand.NewSource(time.Now().Unix()))
	salt := fmt.Sprintf("%06d", rand.Int31())

	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("铲屎官%d", rand.Intn(9000)+1000)

	// 生成初始头像图片
	avatar, _ := imageUtils.ResetAvatar()
	newUser := model.User{
		Username: form.UserName,
		Salt:     salt,
		Password: utils.HandlePassword(form.Password, salt),
		//ClientIp:
		//ClientPort:
		//IsLogout:
		//DeviceId:
		Profile: model.Profile{
			Nickname: name,
			Icon:     avatar,
		},
	}
	// 创建用户
	if result := global.Db.Create(&newUser); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "注册失败:"+result.Error.Error())
		return
	}
	// 验证通过删除验证码信息
	captcha.RedisStore{}.Get(form.CaptchaId, true)
	utils.RespOk(c.Writer, nil, "注册成功")
	return
}
