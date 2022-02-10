package memory

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-kratos/kratos/v2/config"
)

type source struct {
	values map[string]interface{}
}

func New(values map[string]interface{}) config.Source {
	return &source{values: values}
}

func (s *source) Load() ([]*config.KeyValue, error) {
	kvs := make([]*config.KeyValue, 0)
	for key, value := range s.values {
		var result string

		switch reflect.TypeOf(value).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = fmt.Sprintf("%d", value)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = fmt.Sprintf("%d", value)
		case reflect.String:
			result = value.(string)
		case reflect.Bool:
			result = strconv.FormatBool(value.(bool))
		case reflect.Float32, reflect.Float64:
			result = fmt.Sprintf("%f", value)
		}
		kvs = append(kvs, &config.KeyValue{
			Key:   key,
			Value: []byte(result),
		})
	}
	return kvs, nil
}

func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s)
}

type watcher struct {
	source    *source
	closeChan chan struct{}
}

func newWatcher(s *source) (*watcher, error) {
	w := &watcher{
		source:    s,
		closeChan: make(chan struct{}),
	}

	return w, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	<-w.closeChan
	return nil, nil
}

func (w *watcher) Stop() error {
	close(w.closeChan)
	return nil
}
