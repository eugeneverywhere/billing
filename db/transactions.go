package db

import (
	"database/sql"
	"fmt"
	"github.com/eugeneverywhere/billing/db/models"
)

// Transaction - sql transaction interface
type Transaction interface {
	models.Queries

	Commit() error
	Rollback() error
}

type transaction struct {
	connection *sql.Tx

	*models.Managers
	closed bool
}

// Commit - commit the transaction
func (t transaction) Commit() error {
	if t.closed {
		return fmt.Errorf("transaction already closed")
	}

	t.closed = true
	return t.connection.Commit()
}

// Rollback - rollback the transaction
func (t transaction) Rollback() error {
	if t.closed {
		return fmt.Errorf("transaction already closed")
	}

	t.closed = true
	return t.connection.Rollback()
}
