package store

import (
	"database/sql"
	"fem/internal/tokens"
	"time"
)

// types declaration
type PostgresTokenStore struct {
	db *sql.DB
}

// fnc that creates the instance of this type struct --> so methods can use this
func NewPostgresTokenStore (db *sql.DB) *PostgresTokenStore {
	// !intializes the instance of this type struct
	return &PostgresTokenStore{
		db : db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int,ttl time.Duration,scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int,scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int,ttl time.Duration,scope string) (*tokens.Token,error) {
	token,err := tokens.GenerateToken(userID,ttl,scope)
	if err != nil {
		return nil,err
	}
	err = t.Insert(token)
	return token,err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		insert into tokens (hash,user_id,expiry,scope)
		values ($1,$2,$3,$4)
	`
	_,err := t.db.Exec(query,token.Hash,token.UserID,token.Expiry,token.Scope) //* execvte the query n passing placeholder args
		return err

}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userId int,scope string) error {
	query := `
		delete from tokens
		where user_id=$1 and scope=$2
	`

	_,err := t.db.Exec(query,userId,scope)
	return err
}

