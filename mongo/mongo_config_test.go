package mongo

import (
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/liuxiong332/kratos-starter/config/memory"
	"github.com/stretchr/testify/assert"
)

func TestMongoConfig(t *testing.T) {
	cfg := config.New(
		config.WithSource(
			memory.New(map[string]interface{}{
				"username":    "user",
				"password":    "pwd",
				"server":      "server",
				"minPoolSize": 0,
				"maxPoolSize": 1,
			}),
		),
	)
	assert.NoError(t, cfg.Load())

	var mConfig MongoConfig
	assert.NoError(t, cfg.Scan(&mConfig))
	t.Log(mConfig)
}

func TestMongoYamlConfig(t *testing.T) {
	cfg := config.New(
		config.WithSource(
			file.NewSource("./mongo.yaml"),
		),
	)
	assert.NoError(t, cfg.Load())

	var mConfig MongoConfig
	assert.NoError(t, cfg.Value("mongo").Scan(&mConfig))
	t.Log(mConfig)
}
