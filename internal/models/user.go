package models

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type DeleteUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UpdateUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ChangePwdResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
