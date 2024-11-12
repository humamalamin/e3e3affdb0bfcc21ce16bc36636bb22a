package entity

type UserEntity struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

type RequestLogin struct {
	Email    string
	Password string
}
