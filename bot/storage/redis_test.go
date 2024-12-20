package storage

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisStorage(t *testing.T) {
	t.Run("test miniredis", func(t *testing.T) {
		server, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		defer server.Close()

		client := redis.NewClient(&redis.Options{
			Addr: server.Addr(),
		})

		storage := NewRedisStorage(client)

		testStorage(t, storage)
	})

	t.Run("test error handling", func(t *testing.T) {
		client := redis.NewClient(&redis.Options{
			Addr: "invalid.host",
		})

		storage := NewRedisStorage(client)

		var i int
		err := storage.Read("test", "foo", &i)
		assert.Equal(t, "dial tcp: address invalid.host: missing port in address", err.Error())

		keys, err := storage.GetKeys("test")
		assert.Empty(t, keys)
		assert.Equal(t, "dial tcp: address invalid.host: missing port in address", err.Error())
	})
}
