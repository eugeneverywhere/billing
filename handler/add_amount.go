package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eugeneverywhere/billing/types"
)

func (h *handler) handleAddAmount(rawOperation []byte) {
	addAmountData := new(types.AddAmount)
	if err := json.Unmarshal(rawOperation, &addAmountData); err != nil {
		h.log.Errorf("Can't parse add amount operation %q: %v", string(rawOperation), err)
		go h.sendError(&types.OperationResult{
			Operation: &types.Operation{Code: OpAddAmount},
			Result:    ErrWrongFormat,
			Message:   fmt.Sprintf("%v", err),
		})
		return
	}

	h.log.Debugf("Handling: %v", addAmountData)

	err, result := h.addAmount(addAmountData)
	if err != nil || result == nil || result.Result != Ok {
		h.log.Errorf("Adding amount failed for id %v: %v",
			addAmountData.ExternalAccountID, err)
		if result == nil {
			go h.sendError(&types.OperationResult{
				Operation: &types.Operation{Code: addAmountData.Code},
				Result:    ErrInternal,
				Message:   "internal error",
			})
			return
		}
		go h.sendError(result)
	}

}

func (h *handler) addAmount(addAmount *types.AddAmount) (error, *types.OperationResult) {

	account, err := h.db.GetAccountByExternalID(addAmount.ExternalAccountID)
	if err != nil {
		return err, nil
	}

	if account == nil {
		return nil, &types.OperationResult{
			Operation: &types.Operation{Code: addAmount.Code},
			Result:    ErrAccountDoesNotExist,
			Message:   fmt.Sprintf("account %v does not exist", addAmount.ExternalAccountID),
		}
	}

	h.accountMutex.Lock(addAmount.ExternalAccountID)
	defer h.accountMutex.Unlock(addAmount.ExternalAccountID)

	if account.Balance+addAmount.Amount < 0 {
		return nil, &types.OperationResult{
			Operation: &types.Operation{Code: addAmount.Code},
			Result:    ErrInsufficient,
			Message:   fmt.Sprintf("insufficient funds on %v: %v", addAmount.ExternalAccountID, account.Balance),
		}
	}

	account.Balance += addAmount.Amount
	info, err := h.db.UpdateAccountBalance(account)
	if err != nil {
		return err, nil
	}

	return nil, &types.OperationResult{
		Operation: &types.Operation{Code: addAmount.Code},
		Result:    Ok,
		Message:   fmt.Sprintf("now funds on %v: %v", addAmount.ExternalAccountID, info.Balance),
	}
}
