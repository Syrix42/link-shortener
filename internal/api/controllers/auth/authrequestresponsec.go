package auth

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
