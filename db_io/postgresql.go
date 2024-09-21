package db_io

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"database/sql"
	_ "github.com/lib/pq"
)

type PostgresqlDatabase struct {
	Db *sql.DB
	Requests PostgresqlRequests
}

type PostgresqlRequests struct {
	Init   string
	Get    string
	Remove string
	Add    string
}

func PostgresqlRequestsDefault() PostgresqlRequests {
	return PostgresqlRequests {
		Init: `
			CREATE TABLE IF NOT EXISTS public."RefreshTokens"
			(
				access_token_uuid uuid NOT NULL,
				refresh_bcypt bytea NOT NULL,
				PRIMARY KEY (access_token_uuid)
			);
		`,
		Get: `
			SELECT access_token_uuid, refresh_bcypt
			FROM public."RefreshTokens"
			WHERE access_token_uuid = $1;
		`,
		Remove: `
			DELETE FROM public."RefreshTokens"
			WHERE access_token_uuid = $1;
		`,
		Add: `
			INSERT INTO public."RefreshTokens"(access_token_uuid, refresh_bcypt)
			VALUES ($1, $2);
		`,
	}
}

func PostgresqlDatabaseNew(params DatabaseParams, requests PostgresqlRequests) (*PostgresqlDatabase, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", params.Host, params.Port, params.User, params.Password, params.DbName)
         
    conn, err := sql.Open("postgres", psqlconn)
    if err != nil { 
		return nil, errors.New("failed to connect to database: " + err.Error()) 
	}
 
    err = conn.Ping()
    if err != nil { 
		return nil, errors.New("failed to ping database: " + err.Error()) 
	}

    _, err = conn.Exec(requests.Init)
	if err != nil {
		return nil, errors.New("failed to run init request: " + err.Error()) 
	}

	result := PostgresqlDatabase {
		Db: conn,
		Requests: requests,
	}
	return &result, nil 
}

func (psql *PostgresqlDatabase) Close() {
	psql.Db.Close()
}

func (psql *PostgresqlDatabase) Get_RefreshToken(jwtTokenUuid uuid.UUID) (*RefreshToken, error) {
	rows, err := psql.Db.Query(psql.Requests.Get, jwtTokenUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var access_token_uuid uuid.UUID
		var refresh_bcypt []byte

		err = rows.Scan(&access_token_uuid, &refresh_bcypt)
		if err != nil {
			return nil, err
		}

		if access_token_uuid == jwtTokenUuid {
			result := RefreshToken {
				JwtTokenUuid: access_token_uuid,
				RefreshBcryptHash: refresh_bcypt,
			}
			return &result, nil
		}
	}

	return nil, nil
}
func (psql *PostgresqlDatabase) Remove_RefreshToken(jwtTokenUuid uuid.UUID) error {
	_, err := psql.Db.Exec(psql.Requests.Remove, jwtTokenUuid)
	if err != nil {
		return err
	}
	return nil
}
func (psql *PostgresqlDatabase) Add_RefreshToken(token RefreshToken) error {
	_, err := psql.Db.Exec(psql.Requests.Add, token.JwtTokenUuid, token.RefreshBcryptHash)
	if err != nil {
		return err
	}
	return nil
}