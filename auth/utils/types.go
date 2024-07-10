package utils

type BasicUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUserRequest struct {
	Principle string `json:"principle"`
	Password  string `json:"password"`
}