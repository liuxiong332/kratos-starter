package scheduler

import (
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"
	"github.com/RichardKnop/machinery/v2/log"
	"github.com/RichardKnop/machinery/v2/tasks"

	mongobackend "github.com/RichardKnop/machinery/v2/backends/mongo"
	redisbroker "github.com/liuxiong332/kratos-starter/machinery/broker"
	redislock "github.com/liuxiong332/kratos-starter/machinery/lock"

	mongoUtils "github.com/liuxiong332/kratos-starter/mongo"
	redisUtils "github.com/liuxiong332/kratos-starter/redis"
)

type MachineryConfig struct {
	MongoConfig *mongoUtils.MongoConfig
	RedisConfig *redisUtils.RedisConfig
	TaskDBName  string
}

func StartServer(machineryCfg MachineryConfig) (*machinery.Server, error) {
	client, err := mongoUtils.NewClient(machineryCfg.MongoConfig)
	if err != nil {
		return nil, err
	}

	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
		MongoDB: &config.MongoDBConfig{
			Client:   client,
			Database: machineryCfg.TaskDBName,
		},
	}

	redisAddr := machineryCfg.RedisConfig.Nodes
	// Create server instance
	broker := redisbroker.NewGR(cnf, redisAddr, 0, true)
	backend, err := mongobackend.New(cnf)
	if err != nil {
		return nil, err
	}
	lock := redislock.New(cnf, redisAddr, true, 0, 5)
	server := machinery.NewServer(cnf, broker, backend, lock)

	return server, nil
}

func Worker(machineryCfg MachineryConfig) error {
	consumerTag := "machinery_worker"

	server, err := StartServer(machineryCfg)
	if err != nil {
		return err
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorHandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am a start of task handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of task handler for:", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}
