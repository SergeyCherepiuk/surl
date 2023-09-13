package domain

type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
