package store

import (
	"database/sql"
	"fem/internal/tokens"
	"time"
)

//! types declaration
type PostgresTokenStore struct {
	db *sql.DB //* database connection for token operations
}

//! NewPostgresTokenStore --> constructor that creates token store instance
func NewPostgresTokenStore (db *sql.DB) *PostgresTokenStore {
	//* initializes the instance of this type struct
	return &PostgresTokenStore{
		db : db,
	}
}

//! TokenStore interface --> contract for token operations
type TokenStore interface {
	Insert(token *tokens.Token) error //* saves token to database
	CreateNewToken(userID int,ttl time.Duration,scope string) (*tokens.Token, error) //* generates and saves new token
	DeleteAllTokensForUser(userID int,scope string) error //* cleanup old tokens for user
}

//! CreateNewToken --> generates random token and saves it to database
func (t *PostgresTokenStore) CreateNewToken(userID int,ttl time.Duration,scope string) (*tokens.Token,error) {
	//* generate cryptographically secure random token
	token,err := tokens.GenerateToken(userID,ttl,scope)
	if err != nil {
		return nil,err
	}
	//* save token hash to database
	err = t.Insert(token)
	return token,err //* return plaintext token to send to client
}

//! Insert --> saves token hash to database (NOT plaintext for security)
func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		insert into tokens (hash,user_id,expiry,scope)
		values ($1,$2,$3,$4)
	`
	//* execute query with parameterized values (prevents SQL injection)
	_,err := t.db.Exec(query,token.Hash,token.UserID,token.Expiry,token.Scope)
		return err

}

//! DeleteAllTokensForUser --> removes all tokens for specific user and scope
//? useful when user logs out or password changes
func (t *PostgresTokenStore) DeleteAllTokensForUser(userId int,scope string) error {
	query := `
		delete from tokens
		where user_id=$1 and scope=$2
	`

	//* removes all matching tokens from database
	_,err := t.db.Exec(query,userId,scope)
	return err
}

