package main

import (
	"encoding/json"
	"fmt"

	"github.com/liuxiong332/kratos-starter/app"

	"github.com/go-kratos/kratos/v2"
)

func main() {
	appStarter := app.NewApp("ai-studio", nil)

	var kvs map[string]interface{}

	err := appStarter.Config.Scan(&kvs)
	if err != nil {
		fmt.Errorf("Get config error, %v", err)
	}
	if configStr, err := json.MarshalIndent(kvs, "", "  "); err == nil {
		fmt.Printf("Get config: %s", configStr)
	}

	if serverPort, err := appStarter.Config.Value("server.port").String(); err == nil {
		fmt.Printf("Get server port: %s\n", serverPort)
	}

	app := kratos.New(
		kratos.Name("ai-studio"),
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
