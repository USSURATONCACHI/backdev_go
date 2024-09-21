package db_io

import (
	"bytes"
	"errors"
)

type InMemoryDatabase struct {
	Entries []RefreshTokenEntry
}

func (db *InMemoryDatabase) Get_RefreshTokenEntry(entry RefreshTokenEntry) (*RefreshTokenEntry, error) {
	for i := 0; i < len(db.Entries); i++ {
		if bytes.Equal(db.Entries[i].BcryptHash, entry.BcryptHash) {
			result := db.Entries[i].Copy()
			return &result, nil
		}
	}

	return nil, nil
}

func (db *InMemoryDatabase) Add_RefreshTokenEntry(entry RefreshTokenEntry) error {
	for i := 0; i < len(db.Entries); i++ {
		if bytes.Equal(db.Entries[i].BcryptHash, entry.BcryptHash) {
			return errors.New("such refresh token already exists")
		}
	}
	
	db.Entries = append(db.Entries, entry)
	return nil
}

func (db *InMemoryDatabase) Remove_RefreshTokenEntry(entry RefreshTokenEntry) error {
	for i := 0; i < len(db.Entries); i++ {
		if bytes.Equal(db.Entries[i].BcryptHash, entry.BcryptHash) {
			db.Entries = append(db.Entries[:i], db.Entries[i+1:]...)
			return nil
		}
	}
	
	return errors.New("no such refresh token exists")
}