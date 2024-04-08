package ds

import "github.com/google/uuid"

type ProjectAccess struct {
	Id             uuid.UUID `json:"id"`
	ProjectId      uuid.UUID `json:"project_id"`
	CustomerId     uuid.UUID `json:"customer_id"`
	CustomerAccess int       `json:"customer_access"`
}

type AddParticipantReq struct {
	ParticipantEmail string `json:"email"`
	CustomerAccess   string `json:"customer_access"`
}

type UpdateParticipantAccessReq struct {
	ParticipantEmail string `json:"email"`
	CustomerAccess   string `json:"customer_access"`
}

type DeleteParticipantReq struct {
	ParticipantEmail string `json:"email"`
}
