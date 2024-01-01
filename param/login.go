package param

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string   `json:"token"`
	Errors []string `json:"errors"`
}
