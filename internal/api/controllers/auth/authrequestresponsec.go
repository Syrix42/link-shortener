package auth

type RegisterationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterationResponce struct {
	Status   string `json:"status"`
	Messsege string `json:"messege"`
}
