package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"regexp"
)

type User struct {
	gorm.Model
	NickName   string `gorm:"size:32" json:"nick_name"` // 用户昵称
	Username   string `gorm:"size:32" json:"username"`  // 账号
	Password   string `gorm:"size:64" json:"password"`  // 密码
	Avatar     string `gorm:"size:255" json:"avatar"`   // 头像地址
	Tel        string `gorm:"size:18" json:"tel"`       // 电话
	Gender     bool   `json:"gender"`                   // 性别 0女 1男
	Role       int    `json:"role"`                     // 角色 0用户 1管理员
	Status     bool   `json:"status"`                   // 0禁用 1正常
	UserChange UserChange
}

func (u User) getTable() string {
	return "user"
}

type JwtUser struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	Gender   bool   `json:"gender"`
	Tel      string `json:"tel"` // 电话
	Token    string `json:"token"`
}

func ConvertToJwtUser(user User, token string) JwtUser {
	jwtUser := JwtUser{
		UserId:   user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Tel:      user.Tel,
		Token:    token,
	}
	return jwtUser
}

type RequestFrom struct {
	Username string `json:"username" `
	Password string `json:"password" `
	Image    string `json:"imageUtils"`
	IdKey    string `json:"idKey"`
}

// Validate 验证注册请求中的字段
func (r *RequestFrom) Validate() error {
	// 验证用户名长度为4到16位，只包含字母、数字、下划线、减号
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{4,16}$`)
	if !usernameRegex.MatchString(r.Username) {
		return fmt.Errorf("用户名长度为4 ~ 16个字符，只能由字母、数字、“_”和“-”组成")
	}
	// 验证密码长度至少6位，必须包含字母和数字
	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9]{6,16}$`)
	if !passwordRegex.MatchString(r.Password) {
		return fmt.Errorf("密码长度至少6个字符，包含至少一个字母和一个数字")
	}
	// 验证验证码ID不为空
	if r.IdKey == "" {
		return errors.New("需要获取验证码！")
	}
	// 验证验证码不为空且长度在4到6之间
	if r.Image == "" || len(r.Image) < 4 || len(r.Image) > 6 {
		return errors.New("验证码不能为空，且长度必须在4~6个字符之间")
	}
	return nil
}
