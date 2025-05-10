package domain

import "github.com/google/uuid"

type UUIDPagination struct {
	Limit    int
	LastUUID uuid.UUID
}

type IDPagination struct {
	Limit  int
	LastID int
}
