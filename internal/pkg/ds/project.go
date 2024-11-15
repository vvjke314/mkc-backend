package ds

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Id           uuid.UUID `json:"id"`
	OwnerId      uuid.UUID `json:"owner_id"`
	Capacity     int64     `json:"capacity"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
	AdminId      uuid.UUID `json:"admin_id"`
}

type CreateProjectReq struct {
	Name string `json:"name"`
}

type UpdateProjectNameReq struct {
	Name string `json:"name"`
}

type DeleteProjectReq struct {
}

type ProjectData struct {
	Notes []Note `json:"notes"`
	Files []File `json:"files"`
}
