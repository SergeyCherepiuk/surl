package domain

type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type SessionManagerService interface {
	Login(username, password string) error
	Logout(username string) error
}

type AccountManagerService interface {
	Get(username string) (User, error)
	Create(user User) (User, error)
	Update(username string, updates map[string]any) (User, error)
	Delete(user User) (User, error)
}
