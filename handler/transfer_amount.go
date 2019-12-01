package handler

import (
	"fmt"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/eugeneverywhere/billing/types"
)

func (h *handler) TransferAmount(transfer *types.TransferAmount) (error, *types.OperationResult) {
	if transfer.Amount <= 0 {
		return nil, &types.OperationResult{
			Result:  types.ErrNonPositive,
			Message: "amount must be positive",
		}
	}

	err, result, accountSrc := h.checkAccountExists(transfer.Source)
	if accountSrc == nil {
		return err, result
	}

	err, result, accountTgt := h.checkAccountExists(transfer.Target)
	if accountTgt == nil {
		return err, result
	}

	//to avoid deadlocks
	h.accountMutex.Lock(XORStrings(transfer.Source, transfer.Target))
	defer h.accountMutex.Unlock(XORStrings(transfer.Source, transfer.Target))
	//then lock each account to avoid inconsistency
	h.accountMutex.Lock(transfer.Source)
	defer h.accountMutex.Unlock(transfer.Source)
	h.accountMutex.Lock(transfer.Target)
	defer h.accountMutex.Unlock(transfer.Target)

	if accountSrc.Balance < transfer.Amount {
		return nil, &types.OperationResult{
			Result:  types.ErrInsufficient,
			Message: fmt.Sprintf("insufficient on %v", accountSrc.ExternalID),
		}
	}

	accountSrc.Balance -= transfer.Amount
	accountTgt.Balance += transfer.Amount

	tr, err := h.db.Begin()
	if err != nil {
		return err, nil
	}

	_, err = tr.UpdateAccountBalance(accountSrc)
	if err != nil {
		if err = tr.Rollback(); err != nil {
			h.log.Errorf("Rollback failed: %v", err)
		}
		return err, nil
	}

	_, err = tr.UpdateAccountBalance(accountTgt)
	if err != nil {
		if err = tr.Rollback(); err != nil {
			h.log.Errorf("Rollback failed: %v", err)
		}
		return err, nil
	}

	if err := tr.Commit(); err != nil {
		return err, nil
	}

	return nil, &types.OperationResult{
		Result:  types.Ok,
		Message: fmt.Sprintf("transferred %v from %v to %v", transfer.Amount, accountSrc.ExternalID, accountTgt.ExternalID),
	}

}

func (h *handler) checkAccountExists(externalID string) (error, *types.OperationResult, *models.Account) {
	account, err := h.db.GetAccountByExternalID(externalID)
	if err != nil {
		return err, nil, nil
	}

	if account == nil {
		return nil, &types.OperationResult{
			Result:  types.ErrAccountDoesNotExist,
			Message: fmt.Sprintf("account %v does not exist", externalID),
		}, nil
	}
	return nil, nil, account
}
