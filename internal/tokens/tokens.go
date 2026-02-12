package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

//! ScopeAuth --> token type identifier for authentication tokens
const (
	ScopeAuth = "authentication"
)

//! Token struct --> represents authentication token with both plaintext and hashed versions
type Token struct {
	Plaintext string    `json:"token"` //* sent to client (only time they see it)
	Hash      []byte    `json:"-"` //* stored in database for validation
	UserID    int       `json:"-"` //* which user owns this token
	Expiry    time.Time `json:"expiry"` //* when token expires
	Scope     string    `json:"-"` //* token type (authentication, password-reset, etc.)
}

//! GenerateToken --> creates cryptographically secure random token
//! Returns both plaintext (for client) and hash (for database)
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl), //* calculate expiration time
		Scope:  scope,
	}

	//* generate 32 random bytes using crypto/rand (secure random)
	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	//* encode to base32 for URL-safe string (no padding)
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	//* hash plaintext using SHA-256 for database storage
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:] //* convert array to slice
	return token, nil
}