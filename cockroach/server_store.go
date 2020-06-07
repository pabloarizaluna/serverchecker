package cockroach

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pabloarizaluna/serverchecker"
)

type ServerStore struct {
	*sql.DB
}

func (s *ServerStore) Server(id uuid.UUID) (serverchecker.Server, error) {
	var serv serverchecker.Server
	const query = `SELECT * FROM servers WHERE id=$1`
	if err := s.DB.QueryRow(query, id).Scan(
		&serv.DomainID,
		&serv.Address,
		&serv.SslGrade,
		&serv.Country,
		&serv.Owner); err != nil {
		return serverchecker.Server{}, fmt.Errorf("Error getting server: %w", err)
	}
	return serv, nil
}

func (s *ServerStore) Servers(domainID uuid.UUID) ([]serverchecker.Server, error) {
	var ss []serverchecker.Server
	var serv serverchecker.Server
	const query = `SELECT * FROM servers`
	rows, err := s.DB.Query(query)
	if err != nil {
		return []serverchecker.Server{}, fmt.Errorf("Error getting servers: %w", err)
	}
	for rows.Next() {
		if err := rows.Scan(
			&serv.DomainID,
			&serv.Address,
			&serv.SslGrade,
			&serv.Country,
			&serv.Owner); err != nil {
			return []serverchecker.Server{}, fmt.Errorf("Error getting servers: %w", err)
		}
		ss = append(ss, serv)
	}
	return ss, nil
}

func (s *ServerStore) CreateServer(serv *serverchecker.Server) error {
	const query = `INSERT INSTO servers VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := s.DB.Exec(query,
		&serv.ID,
		&serv.DomainID,
		&serv.Address,
		&serv.SslGrade,
		&serv.Country,
		&serv.Owner); err != nil {
		return fmt.Errorf("Error creating server: %w", err)
	}
	return nil
}

func (s *ServerStore) UpdateServer(serv *serverchecker.Server) error {
	const query = `UPDATE servers SET domain_id=$1, address=$2, ssl_grade=$3, country=$4, owner=$5 WHERE id=$6`
	if _, err := s.DB.Exec(query,
		&serv.DomainID,
		&serv.Address,
		&serv.SslGrade,
		&serv.Country,
		&serv.Owner,
		&serv.ID); err != nil {
		return fmt.Errorf("Error updating server: %w", err)
	}
	return nil
}

func (s *ServerStore) DeleteServer(id uuid.UUID) error {
	const query = `DELETE FROM servers WHERE id=$1`
	if _, err := s.DB.Exec(query, id); err != nil {
		return fmt.Errorf("Error deleting server: %w", err)
	}
	return nil
}
