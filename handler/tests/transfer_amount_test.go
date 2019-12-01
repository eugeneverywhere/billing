package tests

import (
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/handler"
	"github.com/eugeneverywhere/billing/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransfer(t *testing.T) {
	mockDB := SetupTests()
	h := handler.NewHandler(mockDB, cache.NewAccountsCache(mockDB))
	testNonPositive(t, h)
	testInsufficientTransfer(t, h)
	testFirstNotExists(t, h)
	testSecondNotExists(t, h)
	testTransferOk(t, h)
}

func testNonPositive(t *testing.T, h handler.Handler) {
	result, _ := h.TransferAmount(&types.TransferAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpTransfer,
		},
		Amount: -1000,
		Source: "Account1",
		Target: "Account2",
	})
	assert.Equal(t, types.ErrNonPositive, result.Result)
}

func testInsufficientTransfer(t *testing.T, h handler.Handler) {
	result, _ := h.TransferAmount(&types.TransferAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpTransfer,
		},
		Amount: 1000000,
		Source: "Account2",
		Target: "Account3",
	})
	assert.Equal(t, types.ErrInsufficient, result.Result)
}

func testFirstNotExists(t *testing.T, h handler.Handler) {
	result, _ := h.TransferAmount(&types.TransferAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpTransfer,
		},
		Amount: 10,
		Source: "Account0",
		Target: "Account2",
	})
	assert.Equal(t, types.ErrAccountDoesNotExist, result.Result)
}

func testSecondNotExists(t *testing.T, h handler.Handler) {
	result, _ := h.TransferAmount(&types.TransferAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpTransfer,
		},
		Amount: 100,
		Source: "Account1",
		Target: "Account0",
	})
	assert.Equal(t, types.ErrAccountDoesNotExist, result.Result)
}

func testTransferOk(t *testing.T, h handler.Handler) {
	result, _ := h.TransferAmount(&types.TransferAmount{
		Operation: &types.Operation{
			ConsumerID:  1,
			OperationID: 1,
			Code:        types.OpTransfer,
		},
		Amount: 100,
		Source: "Account2",
		Target: "Account3",
	})
	assert.Equal(t, types.Ok, result.Result)
}
