package storage

import (
	"database/sql"

	"GitHub.com/sattorovshohruh3009/Authorization/storage/mysql"
	"GitHub.com/sattorovshohruh3009/Authorization/storage/repo"
)

type StorageI interface {
	User() repo.UserStorageI
	Subjects() repo.SubjectsStorageI
	Users_Subjects() repo.Users_SubjectsStorageI
}
type storagePg struct {
	userRepo           repo.UserStorageI
	subjectsRepo       repo.SubjectsStorageI
	users_subjectsRepo repo.Users_SubjectsStorageI
}

func NewStorage(db *sql.DB) StorageI {
	return &storagePg{
		userRepo:           mysql.NewUserStorage(db),
		subjectsRepo:       mysql.NewSubjectsStorage(db),
		users_subjectsRepo: mysql.NewUsers_SubjectsStorage(db),
	}
}
func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}
func (s *storagePg) Subjects() repo.SubjectsStorageI {
	return s.subjectsRepo
}
func (s *storagePg) Users_Subjects() repo.Users_SubjectsStorageI {
	return s.users_subjectsRepo
}
