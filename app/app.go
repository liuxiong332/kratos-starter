package app

import (
	"context"
	"fmt"
	"net/url"
	"os"

	consulConfig "technical-starter/config/consul"

	consulRegistry "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/hashicorp/consul/api"

	vaultConfig "technical-starter/config/vault"

	vaultApi "github.com/hashicorp/vault/api"

	appLog "technical-starter/logger"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

type BootstrapConfig struct {
	ConsulAddress string
	ConsulToken   string
	VaultToken    string
}

func ParseBootstrapConfigEnv() *BootstrapConfig {
	return &BootstrapConfig{
		os.Getenv("APP_CONSUL_ADDRESS"),
		os.Getenv("APP_CONSUL_TOKEN"),
		os.Getenv("APP_VAULT_TOKEN"),
	}
}

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
			u, err := url.Parse(e)
			if err == nil && u.Scheme == "http" {
				return e
			}
		}
	}
	return ""
}

func InitApp(appName string, bootstrapConfig *BootstrapConfig) (log.Logger, registry.Registrar, config.Config) {
	if bootstrapConfig == nil {
		bootstrapConfig = ParseBootstrapConfigEnv()
	}

	// 初始话 logger
	logger := appLog.NewLogger()
	logHelper := log.NewHelper(logger)

	// 初始化 consul config
	client, err := api.NewClient(&api.Config{
		Address: "10.70.2.173:8500",                     // bootstrapConfig.ConsulAddress,
		Token:   "75d7f219-163b-ff19-6730-b5b5b3872e5e", // bootstrapConfig.ConsulToken,
	})
	if err != nil {
		logHelper.Fatal(err)
	}

	consulSrc, err := consulConfig.New(client, consulConfig.WithPath(fmt.Sprintf("config/%s", appName)))

	kvs, err := consulSrc.Load()
	logHelper.Info(kvs)

	// 初始化 consul registry
	registry := consulRegistry.New(client)

	// 初始化 vault config
	vaultAddr := GetInstance(registry, "data-push-provider", logHelper)

	if vaultAddr == "" {
		logHelper.Fatal("Don't find vault instant")
	}

	vaultClient, err := vaultApi.NewClient(&vaultApi.Config{
		Address: vaultAddr,
	})
	if err != nil {
		logHelper.Fatal(err)
	}
	vaultClient.SetToken(bootstrapConfig.VaultToken)

	vaultSrc, err := vaultConfig.New(vaultClient, vaultConfig.WithPath(fmt.Sprintf("config/%s", appName)))

	// 初始化 config
	cfg := config.New(
		config.WithSource(
			consulSrc,
			vaultSrc,
			file.NewSource("./conf/application.yml"),
			env.NewSource("APP_"),
		),
	)
	return logger, registry, cfg
}
