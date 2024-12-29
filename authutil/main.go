package authutil

import (
	"fmt"
	"github.com/Tus1688/library-management-api/types"
	"github.com/essentialkaos/branca/v2"
	"os"
)

type Session interface {
	CreateSessionToken(uid *string) (string, types.Err)
	CreateRefreshToken() (string, types.Err)
	SignRefreshToken(token *string) types.Err
	VerifyRefreshToken(token *string) types.Err
	ValidateToken(token string, ttl uint32) (string, types.Err)
}

type SessionStore struct {
	session       branca.Branca
	refreshSecret string
}

func NewSessionStore() (*SessionStore, error) {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		return nil, fmt.Errorf("SESSION_KEY is not set")
	}
	brc, err := branca.NewBranca([]byte(key))
	if err != nil {
		return nil, err
	}

	refresh := os.Getenv("SESSION_REFRESH_KEY")
	if refresh == "" {
		return nil, fmt.Errorf("SESSION_REFRESH_KEY is not set")
	}

	return &SessionStore{
		session:       brc,
		refreshSecret: refresh,
	}, nil
}
