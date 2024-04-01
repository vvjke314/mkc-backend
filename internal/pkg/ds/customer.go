package ds

import "github.com/google/uuid"

type Customer struct {
	Id         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Type       int       `json:"type"`
}
