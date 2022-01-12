package main

import (
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-zookeeper/zk"
)

type ZkConfig struct {
	Address []string
}

func ParseZkConfig(config config.Config) (*ZkConfig, error) {
	address, err := config.Value("zk.address").String()
	if err != nil {
		return nil, err
	}
	return &ZkConfig{Address: strings.Split(address, ",")}, nil
}

func NewZkClient(zkConfig *ZkConfig) (*zk.Conn, error) {
	zkConn, _, err := zk.Connect(zkConfig.Address, time.Second*5)
	return zkConn, err
}
