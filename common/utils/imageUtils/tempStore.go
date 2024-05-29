package imageUtils

import (
	"context"
	"fmt"
	"golang-pet-api/common/constant"
	"golang-pet-api/common/global"
)

var ctx = context.Background()

type TempRedisStore struct {
}

func (TempRedisStore) Set(id string, value string) error {
	key := constant.Temp_Code + id
	err := global.RedisDb.Set(ctx, key, value, 0).Err()
	return err
}

// Get 查看redis的数据，传入key,和是否删除
func (TempRedisStore) Get(id string, clear bool) string {
	key := constant.Temp_Code + id
	val, err := global.RedisDb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		//clear为true，验证通过，删除这个验证码
		err = global.RedisDb.Del(ctx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}
