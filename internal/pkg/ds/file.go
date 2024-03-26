package ds

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id             uuid.UUID
	ProjectId      uuid.UUID
	Filename       string
	Extension      string
	Size           int
	FilePath       string
	UploadDatetime time.Time
}
