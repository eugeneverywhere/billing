package models

import (
	"database/sql"
	"time"
)

const (
	// DefaultStatus - default items status
	DefaultStatus = 10
	// LockedStatus - locked items status
	LockedStatus = 11
)

// Model - standart db table model
type Model struct {
	ID        int64     `json:"id"`
	Status    int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// SQLConnection - sql connection interface
type SQLConnection interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Queries - queries of all models
type Queries interface {
	AccountQueries
}

// Managers - manager of all models
type Managers struct {
	*AccountManager
}

// NewManager - return new manager instance of all models
func NewManager(connection SQLConnection) *Managers {
	return &Managers{
		&AccountManager{connection},
	}
}
