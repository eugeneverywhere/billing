package tests

import (
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/handler"
	"github.com/eugeneverywhere/billing/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddAmount(t *testing.T) {
	mockDB := SetupTests()
	h := handler.NewHandler(mockDB, cache.NewAccountsCache(mockDB))
	testAccountDoesNotExists(t, h)
	testInsufficient(t, h)
	testAddAmountOk(t, h)
}

func testAccountDoesNotExists(t *testing.T, h handler.Handler) {
	result, _ := h.AddAmount(&types.AddAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpAddAmount,
		},
		Amount:            1000,
		ExternalAccountID: "Account0",
	})
	assert.Equal(t, types.ErrAccountDoesNotExist, result.Result)
}

func testInsufficient(t *testing.T, h handler.Handler) {
	result, _ := h.AddAmount(&types.AddAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpAddAmount,
		},
		Amount:            -200000,
		ExternalAccountID: "Account2",
	})
	assert.Equal(t, types.ErrInsufficient, result.Result)
}

func testAddAmountOk(t *testing.T, h handler.Handler) {
	result, _ := h.AddAmount(&types.AddAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpAddAmount,
		},
		Amount:            100,
		ExternalAccountID: "Account2",
	})
	assert.Equal(t, types.Ok, result.Result)
}
