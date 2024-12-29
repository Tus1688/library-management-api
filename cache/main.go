package cache

import (
	"context"
	"github.com/Tus1688/library-management-api/types"
	"github.com/redis/go-redis/v9"
	"os"
)

type Cache interface {
	Shutdown() error
	SaveRefreshToken(token *string, uid *string) types.Err
}

type RedisStore struct {
	db []*redis.Client
}

func (r *RedisStore) Shutdown() error {
	for _, db := range r.db {
		if err := db.Close(); err != nil {
			return err
		}
	}

	return nil
}

// NewRedisStore creates a new RedisStore instance
// 0 for refresh token
func NewRedisStore(numOfInstance int) (*RedisStore, error) {
	db := make([]*redis.Client, numOfInstance)
	for i := 0; i < numOfInstance; i++ {
		db[i] = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			DB:   i,
		})
		_, err := db[i].Ping(context.TODO()).Result()
		if err != nil {
			return nil, err
		}
	}

	return &RedisStore{
		db: db,
	}, nil
}
