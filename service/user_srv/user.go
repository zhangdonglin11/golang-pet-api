package user_srv

import (
	"golang-pet-api/common/global"
	"golang-pet-api/models/model"
)

type UserService struct{}

// FindUserByLoginFrom 查询用户详情 用户名+密码查询
func (s UserService) FindUserByLoginFrom(from model.RequestFrom) model.User {
	var user model.User
	global.Db.Where("username = ? and password = ?", from.Username, from.Password).Find(&user)
	// 没有错误，说明找到了用户
	return user
}

// FindUserByUserName 用户名查询
func (s UserService) FindUserByUserName(from model.RequestFrom) model.User {
	var user model.User
	global.Db.Where("username = ?", from.Username).Find(&user)
	// 没有错误，说明找到了用户
	return user
}

// 创建用户
func (s UserService) SaveUser(user *model.User) error {
	if err := global.Db.Create(user).Error; err != nil {
		return err
	} else {
		return err
	}
}
