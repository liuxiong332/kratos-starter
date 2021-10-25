package main

import (
	"context"
	"fmt"

	"github.com/liuxiong332/kratos-starter/app"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func customMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		if tr, ok := transport.FromServerContext(ctx); ok {
			fmt.Println("operation:", tr.Operation())
		}
		reply, err = handler(ctx, req)
		return
	}
}

func main() {
	appStarter := app.NewApp("hello-world", nil)

	router := gin.Default()
	// 使用kratos中间件
	router.Use(kgin.Middlewares(recovery.Recovery(), customMiddleware))

	router.GET("/helloworld/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		if name == "error" {
			// 返回kratos error
			kgin.Error(ctx, errors.Unauthorized("auth_error", "no authentication"))
		} else {
			ctx.JSON(200, map[string]string{"welcome": name})
		}
	})

	httpSrv := http.NewServer(http.Address(":8000"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("hello-world"),
		kratos.Version(""),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(appStarter.Logger),
		kratos.Registrar(appStarter.Registry),
		kratos.Server(httpSrv),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
