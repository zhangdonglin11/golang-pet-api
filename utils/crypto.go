package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5Encode(data string) string {
	// 使用 MD5 哈希函数
	hash := md5.New()
	hash.Write([]byte(data))
	hashedPassword := hash.Sum(nil)
	// 将字节转换为十六进制字符串表示
	return hex.EncodeToString(hashedPassword)
}
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

func HandlePassword(password, salt string) string {
	return MD5Encode(password + salt)
}

// VerifyPassword 传未处理密码、盐、处理后的密码
func VerifyPassword(password, salt, hashedPassword string) bool {
	// 比较计算出的哈希值和存储的哈希值是否一致
	return HandlePassword(password, salt) == hashedPassword
}
