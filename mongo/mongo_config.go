package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	configUtil "github.com/liuxiong332/kratos-starter/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Username string
	Password string
	Server   string
}

func ParseMongoConfig(config config.Config) (*MongoConfig, error) {
	var mongoConfig MongoConfig
	configUtil.DumpConfig(config, "mongo", &mongoConfig)
	if mongoConfig.Server == "" {
		return &mongoConfig, fmt.Errorf("Mongo server is not specified")
	}
	return &mongoConfig, nil
}

func NewClient(config *MongoConfig) (*mongo.Client, error) {
	var mongoUri string
	if config.Username != "" {
		mongoUri = fmt.Sprintf("mongodb://%s:%s@%s/?retryWrites=false", config.Username, config.Password, config.Server)
	} else {
		mongoUri = fmt.Sprintf("mongodb://%s/?retryWrites=false", config.Server)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
}
