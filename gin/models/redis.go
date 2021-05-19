package models

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"go-study/gin/conf"
	"go.uber.org/zap"
)

var (
	redisDb *redis.Client
)

func init() {

	redisDb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.C.Redis.IP, conf.C.Redis.Port),
		Password: "",
	})

	_, err = redisDb.Ping().Result()

	if err != nil {
		panic(err)
	}

	zap.L().Debug("redis start")
}

func CloseRedisDb() {

	if redisDb != nil {
		_ = redisDb.Close()
	}

}

func SetString(key string, value interface{}, min int) error {
	return redisDb.Set(key, value, time.Minute*time.Duration(min)).Err()
}

func GetString(key string) (string, error) {
	return redisDb.Get(key).Result()
}

func GetInt(key string) (int, error) {
	return redisDb.Get(key).Int()
}

func DeleteString(key string) (int64, error) {
	return redisDb.Del(key).Result()
}

func IncrString(key string) (int64, error) {
	return redisDb.Incr(key).Result()
}

func Expire(key string, min int) error {
	return redisDb.Expire(key, time.Second*time.Duration(min)).Err()
}
