package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	DB DBTX
}

func NewRepository(DB DBTX) Repository {
	return &repository{
		DB: DB,
	}
}

func (repository *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := repository.DB.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil

}

func (repository *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := repository.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return &User{}, nil
	}

	return &user, nil
}
