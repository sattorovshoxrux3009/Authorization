package models

type Users_Subjects struct {
	Id        int     `json:"id"`
	UserId    int     `json:"user_id"`
	SubjectId int     `json:"subject_id"`
	Score     float64 `json:"score"`
}
type CreateUsers_Subjects struct {
	UserId    int `json:"user_id"`
	SubjectId int `json:"subject_id"`
}
