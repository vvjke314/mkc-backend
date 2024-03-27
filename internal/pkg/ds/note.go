package ds

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id             uuid.UUID
	ProjectId      uuid.UUID
	Title          string
	Content        string
	UpdateDatetime time.Time
	Deadline       time.Time
}
