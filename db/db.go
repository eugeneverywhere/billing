package db

import (
	"database/sql"
	"github.com/eugeneverywhere/billing/db/models"
	// sql driver
	_ "github.com/go-sql-driver/mysql"
)

// DefaultMaxOpenConnection - represents default max open connection to db option
const DefaultMaxOpenConnection = 150

// DB - database interface
type DB interface {
	models.Queries
	Connect() error
	Begin() (Transaction, error)
	Close()
}

type db struct {
	connection         *sql.DB
	connectionURL      string
	maxOpenConnections int

	*models.Managers
}

// New - return new db instance
func New(connectionURL string, maxOpenConnectionsOption ...int) DB {
	maxOpenConnections := DefaultMaxOpenConnection

	if len(maxOpenConnectionsOption) > 0 && maxOpenConnectionsOption[0] != 0 {
		maxOpenConnections = maxOpenConnectionsOption[0]
	}

	return &db{connectionURL: connectionURL, maxOpenConnections: maxOpenConnections}
}

// Connect - connect to db
func (db *db) Connect() error {
	if db.connection != nil {
		db.connection.Close()
	}

	sqlDB, err := sql.Open(MySQLDBDriver, db.connectionURL)
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(db.maxOpenConnections)

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	db.connection = sqlDB
	db.Managers = models.NewManager(db.connection)

	return nil
}

// Begin - begin the transaction
func (db *db) Begin() (Transaction, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{
		connection: tx,
		Managers:   models.NewManager(tx),
	}, nil
}

// Close - close connection to db
func (db *db) Close() {
	db.connection.Close()
}
