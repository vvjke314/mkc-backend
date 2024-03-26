package ds

import "github.com/google/uuid"

type Customer struct {
	Id         uuid.UUID
	FirstName  string
	SecondName string
	Login      string
	Password   string
	Email      string
	Type       int
}
