package db_io

import (
	"github.com/google/uuid"
)

type RefreshToken struct {
	JwtTokenUuid uuid.UUID
	RefreshBcryptHash []byte
}

func (entry RefreshToken) Copy() RefreshToken {
	result := RefreshToken {
		JwtTokenUuid: entry.JwtTokenUuid,
		RefreshBcryptHash: make([]byte, len(entry.RefreshBcryptHash)),
	}

	copy(result.RefreshBcryptHash, entry.RefreshBcryptHash)

	return result
}


type Database interface {
	Get_RefreshToken(jwtTokenUuid uuid.UUID) (*RefreshToken, error);
	Remove_RefreshToken(jwtTokenUuid uuid.UUID) error;
	Add_RefreshToken(token RefreshToken) error;
}