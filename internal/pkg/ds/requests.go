package ds

type SignUpCustomerReq struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Email      string `json:"email"`
}

type LoginCustomerReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
