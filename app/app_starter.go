package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	consulConfig "github.com/liuxiong332/kratos-starter/config/consul"

	"github.com/liuxiong332/kratos-starter/registry/consul"
	consulRegistry "github.com/liuxiong332/kratos-starter/registry/consul"

	"github.com/hashicorp/consul/api"

	vaultConfig "github.com/liuxiong332/kratos-starter/config/vault"

	vaultApi "github.com/hashicorp/vault/api"

	appLog "github.com/liuxiong332/kratos-starter/logger"

	zapLog "github.com/liuxiong332/kratos-starter/logger/zap"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

func GetInstance(discovery registry.Discovery, serviceName string, logHelper *log.Helper) string {
	watcher, err := discovery.Watch(context.Background(), serviceName)
	if err != nil {
		logHelper.Fatal(err)
	}
	svcInstants, err := watcher.Next()
	if err != nil {
		logHelper.Fatal(err)
	}
	if err := watcher.Stop(); err != nil {
		logHelper.Errorf("Failed to http client watch stop, err: %v", err)
	}

	for _, svc := range svcInstants {
		for _, e := range svc.Endpoints {
			return e
		}
	}
	return ""
}

type AppStarter struct {
	Logger   *zapLog.Logger
	Registry *consul.Registry
	Config   config.Config
}

func NewApp(appName string, bootstrapConfig *BootstrapConfig) *AppStarter {
	if bootstrapConfig == nil {
		bootstrapConfig = ParseBootstrapConfigEnv()
	}

	// 初始话 logger
	logger := appLog.NewLogger()
	logHelper := log.NewHelper(logger)

	// 初始化 consul config
	logHelper.Info("Start init consul config")
	client, err := api.NewClient(&api.Config{
		Address: bootstrapConfig.ConsulAddress,
		Token:   bootstrapConfig.ConsulToken,
	})
	if err != nil {
		logHelper.Fatal(err)
	}

	consulSrc, err := consulConfig.New(client, consulConfig.WithPath(fmt.Sprintf("config/%s", appName)))

	// kvs, err := consulSrc.Load()
	// if err != nil {
	// 	logHelper.Fatalf("Load consul config error: %v", err.Error())
	// }
	// logHelper.Info("Get Consul config: %+v", kvs)

	// 初始化 consul registry
	registry := consulRegistry.New(client)

	logHelper.Info("Start init vault config")
	// 初始化 vault config
	vaultAddr := GetInstance(registry, "vault", logHelper)

	if vaultAddr == "" {
		logHelper.Fatal("Don't find vault instant")
	}

	if !strings.HasPrefix(vaultAddr, "http://") {
		vaultAddr = "http://" + vaultAddr
	}

	vaultClient, err := vaultApi.NewClient(&vaultApi.Config{
		Address: vaultAddr,
	})
	if err != nil {
		logHelper.Fatal(err)
	}
	vaultClient.SetToken(bootstrapConfig.VaultToken)

	vaultSrc, err := vaultConfig.New(vaultClient, vaultConfig.WithPath(fmt.Sprintf("secret/%s", appName)))

	// vaultKvs, err := vaultSrc.Load()
	// if err != nil {
	// 	logHelper.Fatalf("Load vault config error: %v", err)
	// }
	// pp.Printf("Get vault kv: %v", vaultKvs)

	// 初始化 config
	configPath := bootstrapConfig.ConfigPath
	if configPath == "" {
		configPath = "./conf/application.yaml"
	}
	var cfg config.Config
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		cfg = config.New(
			config.WithSource(
				consulSrc,
				vaultSrc,
				env.NewSource("APP_"),
			),
		)
	} else {
		cfg = config.New(
			config.WithSource(
				file.NewSource("./conf/application.yaml"),
				consulSrc,
				vaultSrc,
				env.NewSource("APP_"),
			),
		)
	}

	if err := cfg.Load(); err != nil {
		logHelper.Fatalf("App load config error")
	}
	return &AppStarter{
		Logger:   logger,
		Registry: registry,
		Config:   cfg,
	}
}
