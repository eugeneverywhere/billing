package tests

import (
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/handler"
	"github.com/eugeneverywhere/billing/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	mockDB := SetupTests()
	h := handler.NewHandler(mockDB, cache.NewAccountsCache(mockDB))
	testIdTooLong(t, h)
	testIdEmpty(t, h)
	testIdSpaces(t, h)
	testAlreadyExists(t, h)
	testCreateAccountOk(t, h)
}

func testIdTooLong(t *testing.T, h handler.Handler) {
	result, _ := h.CreateAccount(&types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: "VeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongID",
	})
	assert.Equal(t, types.ErrIDTooLong, result.Result)
}

func testIdEmpty(t *testing.T, h handler.Handler) {
	result, _ := h.CreateAccount(&types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: "",
	})
	assert.Equal(t, types.ErrEmptyID, result.Result)
}

func testIdSpaces(t *testing.T, h handler.Handler) {
	result, _ := h.CreateAccount(&types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: "bla bla bla",
	})
	assert.Equal(t, types.ErrSpaces, result.Result)

	result, _ = h.CreateAccount(&types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: "	blabla	",
	})
	assert.Equal(t, types.ErrSpaces, result.Result)
}

func testAlreadyExists(t *testing.T, h handler.Handler) {
	result, _ := h.CreateAccount(&types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: "Account1",
	})
	assert.Equal(t, types.ErrAccountAlreadyExists, result.Result)
}

func testCreateAccountOk(t *testing.T, h handler.Handler) {
	result, _ := h.CreateAccount(&types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: "Account4",
	})
	assert.Equal(t, types.Ok, result.Result)
}
