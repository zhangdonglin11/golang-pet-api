package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unicode"
)

type UserChange struct {
	ID       uint
	User     *User  // 要改成指针，不然就嵌套引用了
	UserID   uint   // 外键_用户id
	NickName string `gorm:"size:32" json:"nick_name"`  // 新昵称
	Avatar   string `gorm:"size:255" json:"avatar"`    // 新图片
	Gender   bool   `json:"gender"`                    // 性别 0女 1男
	Tel      string `gorm:"size:18" json:"tel"`        // 电话
	Approved bool   `gorm:"default:1" json:"approved"` // 0未审核 1已审核
}

func (u *UserChange) MarshalJSON() ([]byte, error) {
	type Alias UserChange
	avatarWithSuffix := u.Avatar + "_suffix"
	return json.Marshal(&struct {
		Avatar string `json:"avatar"`
		*Alias
	}{
		Avatar: avatarWithSuffix,
		Alias:  (*Alias)(u),
	})
}

// 修改的请求结构体
type PutUserFrom struct {
	NickName string `json:"nick_name"`
	Gender   bool   `json:"gender"` // 性别 0女 1男
	Tel      string `json:"tel"`
}

func (p PutUserFrom) Validate() error {
	// 定义允许的字符集合，包括汉字、英文、数字、-和_
	nickNameRegex := regexp.MustCompile("^[\\w\\-\u4e00-\u9fa5]+$")
	if p.NickName != "" && !nickNameRegex.MatchString(p.NickName) {
		return fmt.Errorf("昵称格式汉字、英文、数字、-和_")
	}
	// 统计 名字长度
	hanCount, codeCount := 0, 0
	for _, chat := range p.NickName {
		if unicode.Is(unicode.Han, chat) {
			hanCount += 2
		} else {
			codeCount += 1
		}
	}
	if hanCount > 12 {
		return fmt.Errorf("昵称格式汉字长度最大6个")
	}
	if (hanCount + codeCount) > 12 {
		return fmt.Errorf("昵称格式最多12字节")
	}
	// 电话验证
	telRegex := regexp.MustCompile(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`)
	if p.Tel != "" && !telRegex.MatchString(p.Tel) {
		return fmt.Errorf("请输入正确格式电话号码")
	}
	return nil
}

func UserConvertToUserChange(user User, userChange *UserChange) {
	userChange.UserID = user.ID
	userChange.NickName = user.NickName
	userChange.Avatar = user.Avatar
	userChange.Gender = user.Gender
	userChange.Tel = user.Tel
}
func UserChangeConvertToUser(userChange UserChange, user *User) {
	user.ID = userChange.UserID
	user.NickName = userChange.NickName
	user.Avatar = userChange.Avatar
	user.Gender = userChange.Gender
	user.Tel = userChange.Tel
}
