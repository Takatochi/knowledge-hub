package models

type User struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"johndoe@example.com"`
	Password string `json:"password" example:"password"`
}
