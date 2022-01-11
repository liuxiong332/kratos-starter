package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	data map[string]string
}

func (c testConfig) GetString(key string) (string, error) {
	return c.data[key], nil
}

type Item struct {
	Key1 string
	Key2 string
}

func TestDumpConfig(t *testing.T) {
	var item Item
	config := testConfig{
		data: map[string]string{"ws1.ws2.key1": "value1", "ws2.key1": "12", "ws1.ws2.key2": "value2"},
	}
	dumpConfig(config, "ws1.ws2", &item)
	assert.Equal(t, item, Item{Key1: "value1", Key2: "value2"})
}
