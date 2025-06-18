package model

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Posts    []string `json:"posts"`
}

func NewUser(username string) *User {
	return &User{
		Username: username,
	}
}

func (u *User) GetTypeName() string {
	return "User"
}
