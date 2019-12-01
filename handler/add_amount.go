package handler

import (
	"fmt"
	"github.com/eugeneverywhere/billing/types"
)

func (h *handler) AddAmount(addAmount *types.AddAmount) (*types.OperationResult, error) {
	h.accountMutex.Lock(addAmount.ExternalAccountID)
	defer h.accountMutex.Unlock(addAmount.ExternalAccountID)

	account, err := h.db.GetAccountByExternalID(addAmount.ExternalAccountID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return &types.OperationResult{
			Result:  types.ErrAccountDoesNotExist,
			Message: fmt.Sprintf("account %v does not exist", addAmount.ExternalAccountID),
		}, nil
	}

	if account.Balance+addAmount.Amount < 0 {
		return &types.OperationResult{
			Result:  types.ErrInsufficient,
			Message: fmt.Sprintf("insufficient funds on %v: %v", addAmount.ExternalAccountID, account.Balance),
		}, nil
	}

	account.Balance += addAmount.Amount
	info, err := h.db.UpdateAccountBalance(account)
	if err != nil {
		return nil, err
	}

	return &types.OperationResult{
		Result:  types.Ok,
		Message: fmt.Sprintf("now funds on %v: %v", addAmount.ExternalAccountID, info.Balance),
	}, nil
}
