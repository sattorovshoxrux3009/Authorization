package models

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Auth_time string `json:"auth_time"`
}
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
