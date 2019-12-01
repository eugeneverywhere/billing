package models

import (
	"database/sql"
)

// Account - representation of accounts table
type Account struct {
	Model
	ExternalID string
	Balance    float64
}

// AccountQueries - represents allowed account queries
type AccountQueries interface {
	GetAllAccounts() ([]*Account, error)
	CreateAccount(info *Account) (*Account, error)
	GetAccountByExternalID(externalID string) (*Account, error)
	UpdateAccountBalance(info *Account) (*Account, error)
	LockAccount(id int) error
	UnlockAccount(id int) error
}

// AccountManager - account model manager
type AccountManager struct {
	connection SQLConnection
}

// GetAllAccounts - return all accounts from db
func (m *AccountManager) GetAllAccounts() ([]*Account, error) {
	rows, err := m.connection.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return m.parseMultipleRows(rows)
}

func (m *AccountManager) GetAccountByExternalID(externalID string) (*Account, error) {
	rows, err := m.connection.Query("SELECT * FROM accounts WHERE external_id = ?", externalID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return m.parseSingleRow(rows)
}

// CreateAccount - create new account
func (m *AccountManager) CreateAccount(info *Account) (*Account, error) {
	res, err := m.connection.Exec(`INSERT INTO accounts (external_id, balance, status)
		VALUES (?, ?, ?)`, info.ExternalID, info.Balance, DefaultStatus)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	info.ID = lastID

	return info, err
}

// UpdateAccountBalance - update account balance
func (m *AccountManager) UpdateAccountBalance(info *Account) (*Account, error) {
	_, err := m.connection.Exec(`UPDATE accounts SET balance = ? WHERE id = ?`,
		info.Balance, info.ID)
	if err != nil {
		return nil, err
	}

	return info, err
}

// LockAccount - lock account
func (m *AccountManager) LockAccount(id int) error {
	_, err := m.connection.Exec(`UPDATE accounts SET status = ? WHERE id = ?`, LockedStatus, id)
	return err
}

// UnlockAccount - unlock account
func (m *AccountManager) UnlockAccount(id int) error {
	_, err := m.connection.Exec(`UPDATE accounts SET status = ? WHERE id = ?`, DefaultStatus, id)
	return err
}

func (m *AccountManager) parseMultipleRows(rows *sql.Rows) ([]*Account, error) {
	accounts := make([]*Account, 0)

	for rows.Next() {
		a := new(Account)

		if err := rows.Scan(&a.ID, &a.ExternalID, &a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return accounts, err
		}

		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

func (m AccountManager) parseSingleRow(rows *sql.Rows) (*Account, error) {
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	a := new(Account)

	return a, rows.Scan(&a.ID, &a.ExternalID, &a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt)
}
