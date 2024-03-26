package ds

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Id           uuid.UUID
	OwnerId      uuid.UUID
	Capacity     int
	Name         string
	CreationDate time.Time
	AdminId      uuid.UUID
}
