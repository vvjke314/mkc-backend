package ds

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id             uuid.UUID `json:"id"`
	ProjectId      uuid.UUID `json:"project_id"`
	Filename       string    `json:"filename"`
	Extension      string    `json:"extension"`
	Size           int       `json:"size"`
	FilePath       string    `json:"file_path"`
	UpdateDatetime time.Time `json:"update_datetime"`
}

type DeleteFileReq struct {
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
}

type UpdateFileNameReq struct {
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
}
