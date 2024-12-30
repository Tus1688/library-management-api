package authutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/Tus1688/library-management-api/types"
	"strings"
)

// CreateSessionToken generates a session token for a given user ID.
// It returns the session token as a string and an error if any occurs.
func (s *SessionStore) CreateSessionToken(uid *string) (string, types.Err) {
	brc, err := s.session.EncodeToString([]byte(*uid))
	if err != nil {
		return "", types.Err{Error: "error creating session token"}
	}

	return brc, types.Err{}
}

// CreateRefreshToken generates a signed refresh token.
// It returns the refresh token as a string and an error if any occurs.
func (s *SessionStore) CreateRefreshToken() (string, types.Err) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", types.Err{Error: "error creating refresh token"}
	}
	return base64.URLEncoding.EncodeToString(b), types.Err{}
}

// SignRefreshToken signs the given refresh token using HMAC with SHA-256.
// This function should be called after storing the refresh token in Redis to keep the database clean.
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

// VerifyRefreshToken verifies the given refresh token by checking its signature.
// If the token is valid, it modifies the token to the original token.
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

// ValidateToken validates the given token and returns the user ID if the token is valid.
// It checks if the token is expired based on the provided TTL (time-to-live).
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
