package db_io

type RefreshTokenEntry struct {
	BcryptHash []byte
}

func (entry RefreshTokenEntry) Copy() RefreshTokenEntry {
	return RefreshTokenEntry {
		BcryptHash: entry.BcryptHash,
	}
}


type Database interface {
	Get_RefreshTokenEntry(entry RefreshTokenEntry) (*RefreshTokenEntry, error);
	Add_RefreshTokenEntry(entry RefreshTokenEntry) error;
	Remove_RefreshTokenEntry(entry RefreshTokenEntry) error;
}