package repo

import "context"

type SubjectsStorageI interface {
	Create(ctx context.Context, subject string) (*Subjects, error)
	GetByName(ctx context.Context, name string) ([]Subjects, error)
	GetById(ctx context.Context, id int) (*Subjects, error)
	DeleteByName(ctx context.Context, name string) error
	DeleteById(ctx context.Context, id int) error
}
type Subjects struct {
	Id   int
	Name string
}
