# kratos-starter

Kratos starter is the golang micro service framework starter that based on [kratos](https://github.com/go-kratos/kratos)

It will initialize the following components

### Config

#### config file
use env `APP_CONFIG_PATH`, flag `--config_path` or default `./conf/application.yaml` as the config path

#### consul 
env `APP_CONSUL_ADDRESS` or flag `--consul_address` as the consul address, env `APP_CONSUL_TOKEN` or flag `--consul_token` as the consul token
the config is read from `config/{appName}`, while the app registers itself under `AppStarter.ServiceName`; both are `<appName>` unless env `APP_SERVICE_NAME` or flag `--service_name` sets the registration name separately

#### vault
the vault address service discovery with consul, env `APP_VAULT_TOKEN` or flag `--vault_token` as the vault token)
the config path is `<prefix>/<appName>`, env `APP_VAULT_PREFIX` or flag `--vault_prefix` as the prefix, default `secret` (e.g. set `vega-secret` when the app config lives in the `vega-secret` mount)

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
	appStarter := app.NewApp("my-worker", nil)

	app := kratos.New(
		kratos.Name(appStarter.ServiceName),
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