package models

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	ApiKey   string `json:"apikey"`
}
