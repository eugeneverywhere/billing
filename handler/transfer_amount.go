package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/eugeneverywhere/billing/types"
)

func (h *handler) handleTransfer(rawOperation []byte) {
	transferData := new(types.TransferAmount)
	if err := json.Unmarshal(rawOperation, &transferData); err != nil {
		h.log.Errorf("Can't parse transfer operation %q: %v", string(rawOperation), err)
		go h.sendError(&types.OperationResult{
			Operation: &types.Operation{Code: OpTransfer},
			Result:    ErrWrongFormat,
			Message:   fmt.Sprintf("%v", err),
		})
		return
	}
	h.log.Debugf("Handling: %v", transferData)
	if transferData.Amount <= 0 {
		go h.sendError(&types.OperationResult{
			Operation: &types.Operation{Code: transferData.Code},
			Result:    ErrNonPositive,
			Message:   "amount must be positive",
		})
		return
	}

	err, result := h.transferAmount(transferData)
	if err != nil || result == nil || result.Result != Ok {
		h.log.Errorf("Transfer failed for %v -> %v: %v",
			transferData.Source, transferData.Target, err)
		if result == nil {
			go h.sendError(&types.OperationResult{
				Operation: &types.Operation{Code: transferData.Code},
				Result:    ErrInternal,
				Message:   "internal error",
			})
			return
		}
		go h.sendError(result)
	}
}

func (h *handler) transferAmount(transfer *types.TransferAmount) (error, *types.OperationResult) {
	err, result, accountSrc := h.checkAccountExists(transfer.Source)
	if accountSrc == nil {
		return err, result
	}

	err, result, accountTgt := h.checkAccountExists(transfer.Target)
	if accountTgt == nil {
		return err, result
	}

	h.accountMutex.Lock(XORStrings(transfer.Source, transfer.Target))
	defer h.accountMutex.Unlock(XORStrings(transfer.Source, transfer.Target))

	if accountSrc.Balance < transfer.Amount {
		return nil, &types.OperationResult{
			Operation: &types.Operation{Code: OpTransfer},
			Result:    ErrInsufficient,
			Message:   fmt.Sprintf("insufficient on %v", accountSrc.ExternalID),
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
		Operation: &types.Operation{Code: OpTransfer},
		Result:    Ok,
		Message:   fmt.Sprintf("transferred %v from %v to %v", transfer.Amount, accountSrc.ExternalID, accountTgt.ExternalID),
	}

}

func (h *handler) checkAccountExists(externalID string) (error, *types.OperationResult, *models.Account) {
	account, err := h.db.GetAccountByExternalID(externalID)
	if err != nil {
		return err, nil, nil
	}

	if account == nil {
		return nil, &types.OperationResult{
			Operation: &types.Operation{Code: OpTransfer},
			Result:    ErrAccountDoesNotExist,
			Message:   fmt.Sprintf("account %v does not exist", externalID),
		}, nil
	}
	return nil, nil, account
}
