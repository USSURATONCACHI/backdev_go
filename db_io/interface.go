package db_io

import "github.com/google/uuid"

type RefreshTokenEntry struct {
	RefreshToken uuid.UUID
	RelatedAccessTokenUuid uuid.UUID
}

type Database interface {
	Get_RefreshTokenEntry(token uuid.UUID) (*RefreshTokenEntry, error);
	Add_RefreshTokenEntry(entry RefreshTokenEntry) *error;
	Remove_RefreshTokenEntry(token uuid.UUID) *error;
}