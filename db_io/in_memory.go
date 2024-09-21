package db_io

import (
	"github.com/google/uuid"
	"errors"
)

type InMemoryDatabase struct {
	Entries []RefreshToken
}

func InMemoryDatabaseNew() *InMemoryDatabase {
	return &InMemoryDatabase {
		Entries: make([]RefreshToken, 0),
	}
}

func (db *InMemoryDatabase) Get_RefreshToken(jwtTokenUuid uuid.UUID) (*RefreshToken, error) {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].JwtTokenUuid == jwtTokenUuid {
			result := db.Entries[i].Copy()
			return &result, nil
		}
	}

	return nil, nil
}

func (db *InMemoryDatabase) Add_RefreshToken(token RefreshToken) error {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].JwtTokenUuid == token.JwtTokenUuid {
			return errors.New("such refresh token already exists")
		}
	}
	
	db.Entries = append(db.Entries, token.Copy())
	return nil
}

func (db *InMemoryDatabase) Remove_RefreshToken(jwtTokenUuid uuid.UUID) error {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].JwtTokenUuid == jwtTokenUuid {
			db.Entries = append(db.Entries[:i], db.Entries[i+1:]...)
			return nil
		}
	}
	
	return errors.New("no such refresh token exists")
}