package handler

import (
	"fmt"
	"github.com/eugeneverywhere/billing/types"
)

func (h *handler) AddAmount(addAmount *types.AddAmount) (error, *types.OperationResult) {
	h.accountMutex.Lock(addAmount.ExternalAccountID)
	defer h.accountMutex.Unlock(addAmount.ExternalAccountID)

	account, err := h.db.GetAccountByExternalID(addAmount.ExternalAccountID)
	if err != nil {
		return err, nil
	}

	if account == nil {
		return nil, &types.OperationResult{
			Result:  types.ErrAccountDoesNotExist,
			Message: fmt.Sprintf("account %v does not exist", addAmount.ExternalAccountID),
		}
	}

	if account.Balance+addAmount.Amount < 0 {
		return nil, &types.OperationResult{
			Result:  types.ErrInsufficient,
			Message: fmt.Sprintf("insufficient funds on %v: %v", addAmount.ExternalAccountID, account.Balance),
		}
	}

	account.Balance += addAmount.Amount
	info, err := h.db.UpdateAccountBalance(account)
	if err != nil {
		return err, nil
	}

	return nil, &types.OperationResult{
		Result:  types.Ok,
		Message: fmt.Sprintf("now funds on %v: %v", addAmount.ExternalAccountID, info.Balance),
	}
}
