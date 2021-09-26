package main

import (
	"technical-starter/app"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	logger, registry, cfg := app.InitApp("webull-technical-insight-worker", nil)
	helper := log.NewHelper(logger)
	helper.Info("Init Done")

	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var items map[string]interface{}
	if err := cfg.Scan(&items); err != nil {
		panic(err)
	}

	app := kratos.New(
		kratos.Name("webull-technical-insight-worker"),
		kratos.Version(""),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Registrar(registry),
		kratos.Server(),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
