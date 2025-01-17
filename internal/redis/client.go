package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(address string, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", address, port),
	})

}
