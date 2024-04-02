package mongo

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	ConnectionString string `json:"connectionString"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Server           string `json:"server"`
	MinPoolSize      int    `json:"minPoolSize"`
	MaxPoolSize      int    `json:"maxPoolSize"`
}

func ParseMongoConfig(config config.Config) (*MongoConfig, error) {
	var mongoConfig MongoConfig
	if err := config.Value("mongo").Scan(&mongoConfig); err != nil {
		return nil, err
	}
	return &mongoConfig, nil
}

// connectionString: mongodb+srv://{username}:{password}@{server}/?retryWrites=true&w=1&readPreference=secondaryPreferred
func NewMongoOptions(config *MongoConfig) *options.ClientOptions {
	var mongoUri string
	if config.ConnectionString != "" {
		mongoUri = regexp.MustCompile(`\{\w+\}`).ReplaceAllStringFunc(config.ConnectionString, func(s string) string {
			switch s[1 : len(s)-1] {
			case "username":
				return config.Username
			case "password":
				return config.Password
			case "server":
				return config.Server
			default:
				return s
			}
		})
	} else if config.Username != "" {
		mongoUri = fmt.Sprintf("mongodb://%s:%s@%s/?retryWrites=false", config.Username, config.Password, config.Server)
	} else {
		mongoUri = fmt.Sprintf("mongodb://%s/?retryWrites=false", config.Server)
	}

	mongoOpts := options.Client().ApplyURI(mongoUri)
	if config.MinPoolSize != 0 {
		mongoOpts.SetMinPoolSize(uint64(config.MinPoolSize))
	}
	if config.MaxPoolSize != 0 {
		mongoOpts.SetMaxPoolSize(uint64(config.MaxPoolSize))
	}
	return mongoOpts
}

func NewClient(config *MongoConfig) (*mongo.Client, error) {
	mongoOpts := NewMongoOptions(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return mongo.Connect(ctx, mongoOpts)
}
