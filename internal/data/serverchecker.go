package data

import "github.com/google/uuid"

type Domain struct {
	ID               uuid.UUID `json:"id"`
	ServersChanged   bool      `json:"servers_changed"`
	SslGrade         string    `json:"ssl_grade"`
	PreviousSslGrade string    `json:"previous_ssl_grade"`
	Logo             string    `json:"logo"`
	Title            string    `json:"title"`
	IsDown           bool      `json:"is_down"`
}

type Server struct {
	ID       uuid.UUID `json:"id"`
	DomainID uuid.UUID `json:"domain_id"`
	Address  string    `json:"address"`
	SslGrade string    `json:"ssl_grade"`
	Country  string    `json:"country"`
	Owner    string    `json:"owner"`
}
