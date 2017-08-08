package redis

import "github.com/go-redis/redis"

// Redis: Global redis instance will init when ENV:REDIS IS SET
var Redis *redis.Client

func NewRedis(redisURL string) *redis.Client {
	if redisOpt, err := redis.ParseURL(redisURL); err != nil {
		return nil
	} else {
		return redis.NewClient(redisOpt)
	}
}
