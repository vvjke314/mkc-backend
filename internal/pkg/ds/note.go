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
	UploadDatetime time.Time
	Deadline       time.Time
}
