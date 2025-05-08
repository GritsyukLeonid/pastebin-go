package model

type User struct {
	id       int64
	username string
	posts    []string
}

func NewUser(username string) *User {
	return &User{
		username: username,
	}
}

func (u *User) ID() int64 {
	return u.id
}

func (u *User) SetID(id int64) {
	u.id = id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) AddPost(hash string) {
	u.posts = append(u.posts, hash)
}

func (u *User) Posts() []string {
	return u.posts
}
