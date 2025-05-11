package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"greenlight.shrishail.dev/internal/validator"
)

const (
	ActivationScope = "activation"
)

type Token struct {
	PlainText string
	Hash []byte
	UserID int64
	Expiry time.Time
	Scope string
}

type TokenModel struct {
	DB *sql.DB
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}
	
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be at least 26 characters long")
}

func (m TokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens(hash, user_id, scope, expiry)
		VALUES($1, $2, $3, $4)
	`

	args := []any{token.Hash, token.UserID, token.Scope, token.Expiry}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

/*
DeleteAllForUser()
Deletes all tokens for a specific user and scope.
*/
func (m TokenModel) DeleteAllForUser(userID int64, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id=$1 AND scope=$2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_,err := m.DB.ExecContext(ctx, query, userID, scope)
	return err
}