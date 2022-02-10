package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Server      string `json:"server"`
	MinPoolSize int    `json:"minPoolSize"`
	MaxPoolSize int    `json:"maxPoolSize"`
}

func ParseMongoConfig(config config.Config) (*MongoConfig, error) {
	var mongoConfig MongoConfig
	if err := config.Value("mongo").Scan(&mongoConfig); err != nil {
		return nil, err
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

	mongoOpts := options.Client().ApplyURI(mongoUri)
	if config.MinPoolSize != 0 {
		mongoOpts.SetMinPoolSize(uint64(config.MinPoolSize))
	}
	if config.MaxPoolSize != 0 {
		mongoOpts.SetMaxPoolSize(uint64(config.MaxPoolSize))
	}
	return mongo.Connect(ctx, mongoOpts)
}
