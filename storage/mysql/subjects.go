package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"GitHub.com/sattorovshohruh3009/Authorization/storage/repo"
)

type subjectsRepo struct {
	db *sql.DB
}

func NewSubjectsStorage(db *sql.DB) repo.SubjectsStorageI {
	return &subjectsRepo{
		db: db,
	}
}
func (r *subjectsRepo) Create(ctx context.Context, subject string) (*repo.Subjects, error) {
	query := `
		INSERT INTO Subjects (
			name
		)VALUES (?);
	`
	_, err := r.db.Exec(query, subject)
	if err != nil {
		return nil, err
	}

	query = "SELECT id, name FROM Subjects ORDER BY id DESC LIMIT 1"
	var subj repo.Subjects
	err = r.db.QueryRow(query).Scan(&subj.Id, &subj.Name)
	if err != nil {
		return nil, err
	}
	return &subj, nil
}

func (r *subjectsRepo) GetByName(ctx context.Context, name string) ([]repo.Subjects, error) {
	query := "SELECT id, name FROM Subjects WHERE name = ?"
	rows, err := r.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Yozuvlar yopilishini ta'minlash

	var subjects []repo.Subjects
	for rows.Next() {
		var subj repo.Subjects
		err := rows.Scan(&subj.Id, &subj.Name)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subj)
	}

	// Xatolikni tekshirish
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *subjectsRepo) GetById(ctx context.Context, id int) (*repo.Subjects, error) {
	query := "SELECT id, name FROM Subjects WHERE id = ? LIMIT 1"
	var subj repo.Subjects
	err := r.db.QueryRow(query, id).Scan(&subj.Id, &subj.Name)
	if err != nil {
		return nil, err
	}
	return &subj, nil
}

func (r *subjectsRepo) DeleteById(ctx context.Context, id int) error {
	tsx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM Subjects WHERE id = ?"
	res, err := tsx.Exec(query, id)
	if err != nil {
		tsx.Rollback()
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tsx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tsx.Rollback()
		return fmt.Errorf("no rows affected, id %d not found", id)
	}

	return tsx.Commit()
}

func (r *subjectsRepo) DeleteByName(ctx context.Context, name string) error {
	tsx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM Subjects WHERE name = ?"
	res, err := tsx.Exec(query, name)
	if err != nil {
		tsx.Rollback()
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tsx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tsx.Rollback()
		return fmt.Errorf("no rows affected, id %s not found", name)
	}
	return tsx.Commit()
}
