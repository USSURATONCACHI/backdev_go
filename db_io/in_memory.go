package db_io

import (
	"errors"

	"github.com/google/uuid"
)

type InMemoryDatabase struct {
	Entries []RefreshTokenEntry
}

func (db *InMemoryDatabase) Get_RefreshTokenEntry(token uuid.UUID) (*RefreshTokenEntry, error) {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].RefreshToken == token {

			found := &db.Entries[i]
			result := RefreshTokenEntry {
				RefreshToken: found.RefreshToken,
				RelatedAccessTokenUuid: found.RelatedAccessTokenUuid,
			}

			return &result, nil
		}
	}

	return nil, nil
}

func (db *InMemoryDatabase) Add_RefreshTokenEntry(entry RefreshTokenEntry) *error {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].RefreshToken == entry.RefreshToken {
			result := errors.New("Entry with such primary UUID already exists")
			return &result
		}
	}
	
	db.Entries = append(db.Entries, entry)
	return nil
}

func (db *InMemoryDatabase) Remove_RefreshTokenEntry(token uuid.UUID) *error {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].RefreshToken == token {
			db.Entries = append(db.Entries[:i], db.Entries[i+1:]...)
			return nil
		}
	}
	
	result := errors.New("No such entry exists")
	return &result
}