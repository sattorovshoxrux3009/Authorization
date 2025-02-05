package repo

import "context"

type Users_SubjectsStorageI interface {
	Create(ctx context.Context, req *Users_Subjects) (*Users_Subjects, error)
	GetById(ctx context.Context, id int) (*Users_Subjects, error)
	GetByUserID(ctx context.Context, userId int) ([]Users_Subjects, error)
	GetBySubjectID(ctx context.Context, subjectId int) ([]Users_Subjects, error)
	Update(ctx context.Context, id int, req Users_Subjects) (*Users_Subjects, error)
}
type Users_Subjects struct {
	Id           int
	UserId       int
	SubjectId    int
	Score        float64
	Subject_name string
}
