package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"GitHub.com/sattorovshohruh3009/Authorization/storage/repo"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) repo.UserStorageI {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, req *repo.UserCreate) (*repo.UserCreate, error) {
	query := `
		INSERT INTO Users (
			username,
			password
		)VALUES (?, ?);
	`
	_, err := r.db.Exec(query, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *UserRepo) Get(ctx context.Context, usrname string) (*repo.User, error) {
	query := `
		SELECT 
			id, username,
			password, token, auth_time
		FROM Users WHERE username=?
	`
	var user repo.User
	var token sql.NullString
	var authTime []byte
	err := r.db.QueryRowContext(ctx, query, usrname).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&token,
		&authTime,
	)

	// Agar foydalanuvchi topilmasa
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with username %s not found", usrname)
	}

	// Boshqa xatoliklar
	if err != nil {
		return nil, err
	}

	// token uchun
	if token.Valid {
		user.Token = token.String
	} else {
		user.Token = "NULL" // NULL bo'lsa, bo'sh string
	}

	// auth_time uchun
	if authTime != nil {
		// []byte ni time.Time ga o'zgartirish
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(authTime))
		if err != nil {
			return nil, fmt.Errorf("error parsing auth_time: %v", err)
		}
		user.Auth_time = parsedTime
	} else {
		// NULL bo'lsa, bo'sh vaqt (zero time)
		user.Auth_time = time.Time{}
	}

	return &user, nil
}

func (r *UserRepo) Delete(ctx context.Context, username string) error {
	tsx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tsx.Rollback()
			panic(p)
		}
	}()

	res, err := tsx.Exec("DELETE FROM Users WHERE username=?", username)
	if err != nil {
		_ = tsx.Rollback()
		return err
	}

	data, err := res.RowsAffected()
	if err != nil {
		_ = tsx.Rollback()
		return err
	}

	if data == 0 {
		_ = tsx.Rollback()
		return sql.ErrNoRows
	}

	return tsx.Commit()
}

func (r *UserRepo) UpdatePassword(ctx context.Context, username, newPassword string) (*repo.User, error) {
	query := "UPDATE Users SET password = ? WHERE username = ?"
	_, err := r.db.Exec(query, newPassword, username)
	if err != nil {
		return nil, err
	}

	var user repo.User
	err = r.db.QueryRow("SELECT id, username, password FROM Users WHERE username = ?", username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
