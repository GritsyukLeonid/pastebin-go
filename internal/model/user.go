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

func (u *User) SetID(id int64) {
	u.ID = id
}

func (u *User) AddPost(hash string) {
	u.Posts = append(u.Posts, hash)
}

func (u *User) GetTypeName() string {
	return "User"
}
