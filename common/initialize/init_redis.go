package initialize

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang-pet-api/common/global"
)

func InitRedis() {
	var ctx = context.Background()
	fmt.Println(global.Config.Redis)
	global.RedisDb = redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Address,
		Password: global.Config.Redis.Password, // 没有密码，默认值
		DB:       0,                            // 默认DB 0
	})
	_, err := global.RedisDb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("redis数据库连接成功！")
}
