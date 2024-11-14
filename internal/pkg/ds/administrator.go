package ds

import "github.com/google/uuid"

type Administrator struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type SignUpAdmin struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetCustomerEmailResponse struct {
	Email string `json:"email"`
}
