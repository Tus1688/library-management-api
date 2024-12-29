package cache

import (
	"context"
	"github.com/Tus1688/library-management-api/types"
	"time"
)

func (r *RedisStore) SaveRefreshToken(token *string, uid *string) types.Err {
	err := r.db[0].Set(context.TODO(), *token, *uid, 24*time.Hour).Err()
	if err != nil {
		return types.Err{Error: "unable to save refresh token"}
	}

	return types.Err{}
}
