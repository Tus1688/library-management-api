package authutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/Tus1688/library-management-api/types"
	"strings"
)

func (s *SessionStore) CreateSessionToken(uid *string) (string, types.Err) {
	brc, err := s.session.EncodeToString([]byte(*uid))
	if err != nil {
		return "", types.Err{Error: "error creating session token"}
	}

	return brc, types.Err{}
}

// CreateRefreshToken creates a signed refresh token
func (s *SessionStore) CreateRefreshToken() (string, types.Err) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", types.Err{Error: "error creating refresh token"}
	}
	return base64.URLEncoding.EncodeToString(b), types.Err{}
}

// SignRefreshToken signs the refresh token
// this function should be called after putting the refresh token in the redis that way we keep our database clean
func (s *SessionStore) SignRefreshToken(token *string) types.Err {
	hash := hmac.New(sha256.New, []byte(s.refreshSecret))
	_, err := hash.Write([]byte(*token))
	if err != nil {
		return types.Err{Error: "error signing refresh token"}
	}
	newToken := *token + "." + base64.URLEncoding.EncodeToString(hash.Sum(nil))
	*token = newToken
	return types.Err{}
}

// VerifyRefreshToken verifies the refresh token and modifies the token to the original token
func (s *SessionStore) VerifyRefreshToken(token *string) types.Err {
	parts := strings.Split(*token, ".")
	if len(parts) != 2 {
		return types.Err{Error: "invalid refresh token"}
	}

	// decode the signature
	signature, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return types.Err{Error: "invalid refresh token"}
	}

	// generate the signature and compare it with the signature in the token
	hash := hmac.New(sha256.New, []byte(s.refreshSecret))
	_, err = hash.Write([]byte(parts[0]))
	if err != nil {
		return types.Err{Error: "error verifying refresh token"}
	}

	if !hmac.Equal(signature, hash.Sum(nil)) {
		return types.Err{Error: "invalid refresh token"}
	}
	*token = parts[0]
	return types.Err{}
}

// ValidateToken validates the token and returns the user id if the token is valid
func (s *SessionStore) ValidateToken(token string, ttl uint32) (string, types.Err) {
	decodedString, err := s.session.DecodeString(token)
	if err != nil {
		return "", types.Err{Error: "invalid token"}
	}
	if decodedString.IsExpired(ttl) {
		return "", types.Err{Error: "token expired"}
	}

	return string(decodedString.Payload()), types.Err{}
}
