package repo

import (
	"context"
	"time"
)

type UserStorageI interface {
	Create(ctx context.Context, req *UserCreate) (*User, error)
	Get(ctx context.Context, username string) (*User, error)
	GetToken(ctx context.Context, userId int) (string, error)
	Delete(ctx context.Context, username string) error
	UpdatePassword(ctx context.Context, username, newPassword string) (*User, error)
	UpdateToken(ctx context.Context, username, token string) error
}
type User struct {
	Id        int
	Username  string
	Password  string
	Token     string
	Auth_time time.Time
}

// Post
type UserCreate struct {
	Username string
	Password string
}
