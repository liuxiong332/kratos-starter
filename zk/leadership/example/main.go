package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-zookeeper/zk"
	"github.com/liuxiong332/kratos-starter/zk/leadership"
)

func main() {
	zkCon, _, _ := zk.Connect([]string{"localhost:2181"}, time.Second*5)
	// assert.NoError(t, err)

	leadership.TakeLeader(zkCon, "/test2", log.NewHelper(log.DefaultLogger), func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			fmt.Println("Hello")
			<-time.Tick(time.Second * 5)
		}
	})
}
