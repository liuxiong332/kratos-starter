package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
)

type ConfigGetter interface {
	GetString(key string) (string, error)
}

func dumpConfig(config ConfigGetter, prefix string, target interface{}) {
	targetVal := reflect.Indirect(reflect.ValueOf(target))
	targetType := reflect.TypeOf(target).Elem()
	for i := 0; i < targetVal.NumField(); i += 1 {
		field := targetVal.Field(i)
		fieldName := targetType.Field(i).Name
		fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]

		if field.Kind() == reflect.String {
			value, err := config.GetString(fmt.Sprintf("%s.%s", prefix, fieldName))
			if err == nil {
				field.SetString(value)
			}
		}
	}
}

type configGetAdapter struct {
	config config.Config
}

func (c configGetAdapter) GetString(key string) (string, error) {
	return c.config.Value(key).String()
}

func DumpConfig(config config.Config, prefix string, target interface{}) {
	dumpConfig(configGetAdapter{config: config}, prefix, target)
}
