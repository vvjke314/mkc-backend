package ds

import "github.com/google/uuid"

type ProjectAccess struct {
	Id             uuid.UUID
	ProjectId      uuid.UUID
	CustomerId     uuid.UUID
	CustomerAccess int
}

type AddParticipantReq struct {
}

type DeleteParticipantReq struct {
}
