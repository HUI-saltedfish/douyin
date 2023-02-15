package redisService

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.ClusterClient
var Ctx = context.Background()

func init() {
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:6374", "localhost:6375", "localhost:6376"},
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
