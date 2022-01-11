package redis

import (
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Nodes []string
}

func ParseRedisConfig(config config.Config) (*RedisConfig, error) {
	nodes, err := config.Value("redis.nodes").String()
	if err != nil {
		return nil, err
	}
	return &RedisConfig{Nodes: strings.Split(nodes, ",")}, nil
}

func NewClient(config *RedisConfig) redis.UniversalClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: config.Nodes,
	})
}
