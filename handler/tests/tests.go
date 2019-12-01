package tests

import (
	"github.com/eugeneverywhere/billing/db/mocks"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/stretchr/testify/mock"
	"time"
)

var allAccounts = []*models.Account{&models.Account{
	Model: models.Model{
		ID:        1,
		Status:    10,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	ExternalID: "Account1",
	Balance:    0,
},
	&models.Account{
		Model: models.Model{
			ID:        2,
			Status:    10,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		ExternalID: "Account2",
		Balance:    100042.31,
	},
	&models.Account{
		Model: models.Model{
			ID:        3,
			Status:    10,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		ExternalID: "Account3",
		Balance:    500.32,
	},
}

func SetupTests() *mocks.DB {
	mockDB := &mocks.DB{}

	mockGetAllAccounts := func() ([]*models.Account, error) {
		return allAccounts, nil
	}

	mockGetAccountByExternalID := func(externalID string) *models.Account {
		switch externalID {
		case "Account1":
			return allAccounts[0]
		case "Account2":
			return allAccounts[1]
		case "Account3":
			return allAccounts[2]
		default:
			return nil
		}
	}

	mockDB.On("GetAllAccounts").Return(mockGetAllAccounts)
	mockDB.On("GetAccountByExternalID",
		mock.AnythingOfType("string")).Return(mockGetAccountByExternalID, nil)
	mockDB.On("CreateAccount",
		mock.AnythingOfType("*models.Account")).Return(func(input *models.Account) *models.Account {
		return input
	}, nil)
	mockDB.On("UpdateAccountBalance",
		mock.AnythingOfType("*models.Account")).Return(func(input *models.Account) *models.Account {
		return input
	}, nil)

	mockTransaction := &mocks.Transaction{}
	mockTransaction.On("Commit").Return(nil)
	mockTransaction.On("UpdateAccountBalance",
		mock.AnythingOfType("*models.Account")).Return(func(input *models.Account) *models.Account {
		return input
	}, nil)

	mockDB.On("Begin").Return(mockTransaction, nil)

	return mockDB
}
