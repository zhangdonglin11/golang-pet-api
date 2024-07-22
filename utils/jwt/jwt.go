package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

// token 过期时间
const TokenExpirationTime = 60 * 60 * 24 * 20
const TokenKey = "SLDSLDKJSSK"

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GetJwtToken(secretKey string, iat, seconds int64, data map[string]interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range data {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// 解析 JWT Token 的函数
func ParseJwtToken(secretKey string, tokenString string) (jwt.MapClaims, error) {
	// 解析 JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("未知的签名方法: %v", token.Header["alg"])
		}
		// 返回用于验证签名的 key，即你在生成 JWT 时使用的 secretKey
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("JWT 解析失败: %v", err)
	}

	// 检查 token 是否有效
	if !token.Valid {
		return nil, fmt.Errorf("无效的JWT token")
	}

	// 提取并返回有效载荷中的数据
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("无效的JWT载荷")
	}
}
