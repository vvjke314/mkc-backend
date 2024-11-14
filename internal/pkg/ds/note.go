package ds

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id             uuid.UUID `json:"id"`
	ProjectId      uuid.UUID `json:"project_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	UpdateDatetime time.Time `json:"update_datetime"`
	Deadline       time.Time `json:"deadline"`
	Overdue        int       `json:"overdue"`
}

type CreateNoteReq struct {
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Deadline time.Time `json:"deadline"`
}

type UpdateNoteDeadlineReq struct {
	Deadline time.Time `json:"deadline"`
}

type DeleteNoteReq struct {
}
