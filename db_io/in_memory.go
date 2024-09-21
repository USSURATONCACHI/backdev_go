package db_io

import (
	"errors"
)

type InMemoryDatabase struct {
	Entries []RefreshTokenEntry
}

func (db *InMemoryDatabase) Get_RefreshTokenEntry(entry RefreshTokenEntry) (*RefreshTokenEntry, error) {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].BcryptHash == entry.BcryptHash {
			result := db.Entries[i].Copy()
			return &result, nil
		}
	}

	return nil, nil
}

func (db *InMemoryDatabase) Add_RefreshTokenEntry(entry RefreshTokenEntry) *error {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].BcryptHash == entry.BcryptHash {
			result := errors.New("Such refresh token already exists")
			return &result
		}
	}
	
	db.Entries = append(db.Entries, entry)
	return nil
}

func (db *InMemoryDatabase) Remove_RefreshTokenEntry(entry RefreshTokenEntry) *error {
	for i := 0; i < len(db.Entries); i++ {
		if db.Entries[i].BcryptHash == entry.BcryptHash {
			db.Entries = append(db.Entries[:i], db.Entries[i+1:]...)
			return nil
		}
	}
	
	result := errors.New("No such refresh token exists")
	return &result
}