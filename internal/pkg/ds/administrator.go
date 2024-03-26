package ds

import "github.com/google/uuid"

type Administrator struct {
	Id       uuid.UUID
	Name     string
	Email    string
	Password string
}
