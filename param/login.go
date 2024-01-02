package param

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	Token  string   `json:"token"`
	Errors []string `json:"errors"`
}
