package redis

import (
	"github.com/go-redis/redis"
	// "github.com/go-redis/redis/v8"
	"fmt"
	"github.com/alicebob/miniredis/v2"
)

var (
	Redis *redis.Client
)

func init() {
	// addr := viper.GetString("redis.addr")
	// password := viper.GetString("redis.password")
	// db := viper.GetInt("redis.db")

	// Redis = redis.NewClient(&redis.Options{
	// 	Addr:     addr,
	// 	Password: password, // no password set
	// 	DB:       db,       // use default DB
	// })

	// pong, err := Redis.Ping().Result()
	// if err != nil {
	// 	logrus.Info("连接redis失败, ", err)
	// 	panic("连接redis失败, redis地址=" + addr + ", error: " + err.Error())
	// }

	// logrus.Info("连接redis成功, 地址=", addr, ", pong=", pong)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// 使用 go-redis 客户端连接到这个内存中的 Redis 服务器
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
		// Password: s.Password(), // 如果设置了密码，则在这里输入
		DB: 0, // 默认数据库
	})

	// 设置一个键值对
	err = rdb.Set("mykey", "value", 0).Err()
	if err != nil {
		fmt.Printf("Error setting key: %v\n", err)
		return
	}

	// 获取一个键的值
	val, err := rdb.Get("mykey").Result()
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
		return
	}

	Redis = rdb

	fmt.Printf("Value for mykey: %s\n", val)
}
