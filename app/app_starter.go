package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	watcher, err := discovery.Watch(ctx, serviceName)
	if err != nil {
		logHelper.Fatal(err)
	}
	svcInstants, err := watcher.Next()
	if err != nil {
		logHelper.Error(err)
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

func newVaultConfig(registry *consul.Registry, logHelper *log.Helper, appName string, bootstrapConfig *BootstrapConfig) config.Source {
	// 初始化 vault config
	vaultAddr := GetInstance(registry, "vault", logHelper)

	if vaultAddr != "" {
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
		if err != nil {
			logHelper.Fatalf("New vault config error: %v", err)
		}
		return vaultSrc
	}
	return nil
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
	if err != nil {
		logHelper.Fatalf("New consul error: %v", err.Error())
	}

	// 初始化 consul registry
	var registry *consul.Registry
	if bootstrapConfig.ConsulTags != "" {
		registry = consulRegistry.New(client, consul.WithTags(strings.Split(bootstrapConfig.ConsulTags, ",")))
	} else {
		registry = consulRegistry.New(client)
	}

	logHelper.Info("Start init vault config")

	vaultSrc := newVaultConfig(registry, logHelper, appName, bootstrapConfig)

	// 初始化 config
	configPath := bootstrapConfig.ConfigPath
	if configPath == "" {
		configPath = "./conf/application.yaml"
	}

	var configSrcs []config.Source

	var cfg config.Config
	if _, err := os.Stat(configPath); !errors.Is(err, os.ErrNotExist) {
		configSrcs = append(configSrcs, file.NewSource("./conf/application.yaml"))
	}

	configSrcs = append(configSrcs, consulSrc)

	if vaultSrc != nil {
		configSrcs = append(configSrcs, vaultSrc)
	}

	configSrcs = append(configSrcs, env.NewSource("APP_"))

	cfg = config.New(config.WithSource(configSrcs...))

	if err := cfg.Load(); err != nil {
		logHelper.Fatalf("App load config error")
	}
	return &AppStarter{
		Logger:   logger,
		Registry: registry,
		Config:   cfg,
	}
}
