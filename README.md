# github.com/liuxiong332/kratos-starter

Kratos starter is the golang micro service framework starter that based on [kratos](https://github.com/go-kratos/kratos)

It will initialize the following components

### Config

#### config file
use env `APP_CONFIG_PATH`, flag `--config_path` or default `./conf/application.yaml` as the config path

#### consul 
env `APP_CONSUL_ADDRESS` or flag `--consul_address` as the consul address, env `APP_CONSUL_TOKEN` or flag `--consul_token` as the consul token

#### vault
the vault address service discovery with consul, env `APP_VAULT_TOKEN` or flag `--vault_token` as the vault token)

#### env
The environment variable with prefix `APP_` will used as the config.

### Log

Initialize zap log library with structure log.

### Registry

Initialize the consul registry.

# Quick start

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/liuxiong332/kratos-starter/app"

	"github.com/go-kratos/kratos/v2"
)

func main() {
	appStarter := app.NewApp("webull-technical-insight-worker", nil)

	app := kratos.New(
		kratos.Name("webull-technical-insight-worker"),
		kratos.Version(""),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(appStarter.Logger),
		kratos.Registrar(appStarter.Registry),
		kratos.Server(),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

```