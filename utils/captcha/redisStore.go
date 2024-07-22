package captcha

import (
	"context"
	"fmt"
	"golang-pet-api/global"
	"time"
)

// redis存取验证码
var ctx = context.Background()

// RedisStore 结构用于操作
type RedisStore struct {
}

// 实现设置captcha的方法
func (r RedisStore) Set(id string, value string) error {
	key := global.LOGIN_CODE + id
	//time.Minute*2：有效时间2分钟
	err := global.RedisDb.Set(ctx, key, value, time.Minute*2).Err()
	return err
}

// 实现获取captcha的方法
func (r RedisStore) Get(id string, clear bool) string {
	key := global.LOGIN_CODE + id
	val, err := global.RedisDb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		//clear为true，验证通过，删除这个验证码
		err := global.RedisDb.Del(ctx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// 实现验证captcha的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := r.Get(id, clear)
	fmt.Println("key:" + id + ";value:" + v + ";answer:" + answer)
	return v == answer
}
