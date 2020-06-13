package cockroach

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	*DomainStore
}

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error connectiong to database: %w", err)
	}

	return &Store{
		DomainStore: &DomainStore{DB: db},
	}, nil
}
