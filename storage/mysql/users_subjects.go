package mysql

import (
	"context"
	"database/sql"

	"GitHub.com/sattorovshohruh3009/Authorization/storage/repo"
)

type users_subjectsRepo struct {
	db *sql.DB
}

func NewUsers_SubjectsStorage(db *sql.DB) repo.Users_SubjectsStorageI {
	return &users_subjectsRepo{
		db: db,
	}
}

func (r *users_subjectsRepo) Create(ctx context.Context, req *repo.Users_Subjects) (*repo.Users_Subjects, error) {
	query := `
		INSERT INTO Users_subjects (
			user_id,
			subject_id
		) VALUES (?, ?);
	`
	_, err := r.db.Exec(query, req.UserId, req.SubjectId)
	if err != nil {
		return nil, err
	}

	query = "SELECT id, user_id, subject_id FROM Users_subjects ORDER BY id DESC LIMIT 1"
	var user_subject repo.Users_Subjects
	err = r.db.QueryRow(query).Scan(&user_subject.Id, &user_subject.UserId, &user_subject.SubjectId)
	if err != nil {
		return nil, err
	}
	return &user_subject, nil
}

func (r *users_subjectsRepo) GetById(ctx context.Context, id int) (*repo.Users_Subjects, error) {
	query := "SELECT id, user_id, subject_id, score FROM Users_subjects WHERE id = ? LIMIT 1"
	var subj repo.Users_Subjects
	var score sql.NullFloat64

	err := r.db.QueryRow(query, id).Scan(&subj.Id, &subj.UserId, &subj.SubjectId, &score)
	if err != nil {
		return nil, err
	}

	if score.Valid {
		subj.Score = score.Float64
	} else {
		subj.Score = 0
	}

	return &subj, nil
}

func (r *users_subjectsRepo) GetByUserID(ctx context.Context, userID int) ([]repo.Users_Subjects, error) {
	query := `SELECT us.id, us.user_id, us.subject_id, COALESCE(us.score, 0) as score, s.name as subject_name 
	          FROM Users_subjects us 
	          JOIN Subjects s ON us.subject_id = s.id 
	          WHERE us.user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []repo.Users_Subjects
	for rows.Next() {
		var subject repo.Users_Subjects
		if err := rows.Scan(&subject.Id, &subject.UserId, &subject.SubjectId, &subject.Score, &subject.Subject_name); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *users_subjectsRepo) GetBySubjectID(ctx context.Context, subjectId int) ([]repo.Users_Subjects, error) {
	query := "SELECT id, user_id, subject_id, score FROM Users_subjects WHERE subject_id = ?"
	rows, err := r.db.Query(query, subjectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []repo.Users_Subjects

	for rows.Next() {
		var subj repo.Users_Subjects
		var score sql.NullFloat64

		err := rows.Scan(&subj.Id, &subj.UserId, &subj.SubjectId, &score)
		if err != nil {
			return nil, err
		}

		if score.Valid {
			subj.Score = score.Float64
		} else {
			subj.Score = 0
		}

		subjects = append(subjects, subj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *users_subjectsRepo) Update(ctx context.Context, id int, req repo.Users_Subjects) (*repo.Users_Subjects, error) {
	query := `
		UPDATE Users_subjects
		SET user_id = ?, subject_id = ?, score = ?
		WHERE id = ?`

	// SQL so'rovini bajarish
	_, err := r.db.Exec(query, req.UserId, req.SubjectId, req.Score, id)
	if err != nil {
		return nil, err
	}

	// Yangilangan yozuvni olish
	var updatedSubj repo.Users_Subjects
	err = r.db.QueryRow("SELECT id, user_id, subject_id, score FROM Users_subjects WHERE id = ?", id).Scan(&updatedSubj.Id, &updatedSubj.UserId, &updatedSubj.SubjectId, &updatedSubj.Score)
	if err != nil {
		return nil, err
	}

	return &updatedSubj, nil
}
