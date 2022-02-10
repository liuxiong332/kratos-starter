package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	src := New(map[string]interface{}{"int": 1, "int8": int8(1), "uint": uint(1), "float": 1.2, "str": "string", "bool": true})
	_, err := src.Load()
	assert.NoError(t, err)

	// assert.Equal(t, kvs[0].Value, []byte("1"))
	// assert.Equal(t, kvs[1].Value, []byte("1"))
	// assert.Equal(t, kvs[2].Value, []byte("1"))
	// assert.Equal(t, kvs[3].Value, []byte("1.200000"))
	// assert.Equal(t, kvs[4].Value, []byte("string"))
	// assert.Equal(t, kvs[5].Value, []byte("true"))
}
