package entity

type User struct {
	ID       string
	Username string
	Password string
	Email    string
}

type UserLoginData struct {
	ID       string
	Username string
	Email    string
}
