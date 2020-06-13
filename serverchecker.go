package serverchecker

import "github.com/google/uuid"

type Domain struct {
	ID               uuid.UUID `json:"id"`
	ServersChanged   bool      `json:"servers_changed"`
	SslGrade         string    `json:"ssl_grade"`
	PreviousSslGrade string    `json:"previous_ssl_grade"`
	Logo             string    `json:"logo"`
	Title            string    `json:"title"`
	IsDown           bool      `json:"is_down"`
	Host             string    `json:"host"`
	Servers          []Server  `json:"servers"`
}

type Server struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

type DomainStore interface {
	Domain(host string) (Domain, error)
	Domains() ([]Domain, error)
	CreateDomain(d *Domain) error
	UpdateDomain(d *Domain) error
	DeleteDomain(id uuid.UUID) error
}

type Store interface {
	DomainStore
}
