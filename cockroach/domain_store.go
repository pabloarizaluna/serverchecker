package cockroach

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pabloarizaluna/serverchecker"
)

type DomainStore struct {
	*sql.DB
}

func (s *DomainStore) Domain(name string) (serverchecker.Domain, error) {
	var d serverchecker.Domain
	const query = `SELECT id, servers_changed, is_down, ssl_grade, previous_ssl_grade, logo, title, name FROM domains WHERE name=$1`
	if err := s.DB.QueryRow(query, name).Scan(
		&d.ID,
		&d.ServersChanged,
		&d.IsDown,
		&d.SslGrade,
		&d.PreviousSslGrade,
		&d.Logo,
		&d.Title,
		&d.Host); err != nil {
		return serverchecker.Domain{}, err
	}

	return d, nil
}

func (s *DomainStore) Domains() ([]serverchecker.Domain, error) {
	var dd []serverchecker.Domain
	var d serverchecker.Domain
	const query = `SELECT id, servers_changed, is_down, ssl_grade, previous_ssl_grade, logo, title, name FROM domains`
	rows, err := s.DB.Query(query)
	if err != nil {
		return []serverchecker.Domain{}, fmt.Errorf("Error getting domains: %w", err)
	}
	for rows.Next() {
		if err := rows.Scan(
			&d.ID,
			&d.ServersChanged,
			&d.IsDown,
			&d.SslGrade,
			&d.PreviousSslGrade,
			&d.Logo,
			&d.Title,
			&d.Host); err != nil {
			return []serverchecker.Domain{}, fmt.Errorf("Error getting domain: %w", err)
		}
		dd = append(dd, d)
	}

	return dd, nil

}

func (s *DomainStore) CreateDomain(d *serverchecker.Domain) error {
	const query = `INSERT INTO domains VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := s.DB.Exec(query,
		d.ID,
		d.ServersChanged,
		d.SslGrade,
		d.PreviousSslGrade,
		d.Logo,
		d.Title,
		d.IsDown,
		d.Host); err != nil {
		return fmt.Errorf("Error creating domain: %w", err)
	}
	return nil
}

func (s *DomainStore) UpdateDomain(d *serverchecker.Domain) error {
	const query = `UPDATE domains SET servers_changed=$1, ssl_grade=$2, previous_ssl_grade=$3, logo=$4, title=$5, is_down=$6, name=$7 WHERE id=$8`
	if _, err := s.DB.Exec(query,
		d.ServersChanged,
		d.SslGrade,
		d.PreviousSslGrade,
		d.Logo,
		d.Title,
		d.IsDown,
		d.Host,
		d.ID); err != nil {
		return fmt.Errorf("Error updating domain: %w", err)
	}
	return nil
}

func (s *DomainStore) DeleteDomain(id uuid.UUID) error {
	const query = `DELETE FROM domains WHERE id=$1`
	if _, err := s.DB.Exec(query, id); err != nil {
		return fmt.Errorf("Error deleting domain: %w", err)
	}
	return nil
}
