package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/constant"
	"golang-pet-api/models/model"
	"time"
)

type userStaClaims struct {
	model.JwtUser
	jwt.StandardClaims
}

// token 过期时间
const TokenExpireDuration = time.Hour * 360

// token 加密
var jwtKey = []byte("golang-pet")
var (
	ErrAbsent  = "token absent"  //令牌不存在
	ErrInvalid = "token invalid" //令牌无效
)

// 根据用户信息生成token
func CreateTokenByUser(user model.User) (string, error) {
	var jwtUser = model.JwtUser{
		UserId:   user.ID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Tel:      user.Tel,
	}
	claims := userStaClaims{
		jwtUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 设置过期时间为 24 小时后
			IssuedAt:  time.Now().Unix(),                          // 设置签发时间为当前时间
			Issuer:    "admin",                                    // 设置签发用户
		}}
	// 创建 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名 token 并获取完整的 token 字符串
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析 JWT
func ValidateJWT(tokenString string) (*model.JwtUser, error) {
	if tokenString == "" {
		return nil, errors.New(ErrAbsent)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法是否匹配
		return jwtKey, nil
	})
	if token == nil {
		return nil, errors.New(ErrInvalid)
	}
	claims := userStaClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return &claims.JwtUser, err
}

// 返回用户id
func GetUserId(c *gin.Context) (uint, error) {
	u, exist := c.Get(constant.ContextKeyUserObj)
	if !exist {
		return 0, errors.New("can't get user id")
	}
	user, ok := u.(*model.JwtUser)
	if ok {
		return user.UserId, nil
	}
	return 0, errors.New("can't convert to id struct")
}

// 返回用户名
func GetNickName(c *gin.Context) (string, error) {
	u, exist := c.Get(constant.ContextKeyUserObj)
	if !exist {
		return string(string(0)), errors.New("can't get user name")
	}
	user, ok := u.(*model.JwtUser)
	if ok {
		return user.NickName, nil
	}
	return string(string(0)), errors.New("can't convert to api name")
}

// 返回用户信息
func GetAdmin(c *gin.Context) (*model.JwtUser, error) {
	u, exist := c.Get(constant.ContextKeyUserObj)
	if !exist {
		return nil, errors.New("can't get api")
	}
	user, ok := u.(*model.JwtUser)
	if ok {
		return user, nil
	}
	return nil, errors.New("can't convert to api struct")
}
