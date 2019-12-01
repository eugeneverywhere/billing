package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/eugeneverywhere/billing/types"
)

func (h *handler) handleAccountCreate(rawOperation []byte) {
	createAccountData := new(types.CreateAccount)
	if err := json.Unmarshal(rawOperation, &createAccountData); err != nil {
		h.log.Errorf("Can't parse create account operation %q: %v", string(rawOperation), err)
		go h.sendError(&types.OperationResult{
			Operation: createAccountData.Operation,
			Result:    ErrWrongFormat,
			Message:   fmt.Sprintf("%v", err),
		})
		return
	}

	h.log.Debugf("Handling: %v", createAccountData)

	err, result := h.CreateAccount(createAccountData)

	if err != nil || result == nil || result.Result != Ok {
		h.log.Errorf("Account creation failed for id %v: %v",
			createAccountData.ExternalAccountID, err)
		if result == nil {
			go h.sendError(&types.OperationResult{
				Operation: createAccountData.Operation,
				Result:    ErrInternal,
				Message:   "internal error",
			})
			return
		}
		result.Operation = createAccountData.Operation
		go h.sendError(result)
	}
}

func (h *handler) accountExists(externalID string) (bool, error) {
	accountsMap := h.accountsCache.GetAccountsByExtID()
	_, ok := accountsMap[externalID]
	if ok {
		return true, nil
	}

	account, err := h.db.GetAccountByExternalID(externalID)
	if err != nil {
		h.log.Errorf("Accessing db failed: :v", err)
		return false, err
	}
	if account != nil {
		return true, nil
	}
	return false, nil
}

func (h *handler) getAccountByExternalID(externalID string) (bool, error) {
	account, err := h.db.GetAccountByExternalID(externalID)
	if err != nil {
		h.log.Errorf("Accessing db failed: :v", err)
		return false, err
	}
	if account != nil {
		return true, nil
	}
	return false, nil
}

func (h *handler) CreateAccount(operation *types.CreateAccount) (error, *types.OperationResult) {
	if ContainsSpaces(operation.ExternalAccountID) {
		return nil, &types.OperationResult{
			Result:  ErrSpaces,
			Message: fmt.Sprintf("blank characters not allowed: %v", operation.ExternalAccountID),
		}
	}

	if operation.ExternalAccountID == "" {
		return nil, &types.OperationResult{
			Result:  ErrEmptyID,
			Message: "empty id not allowed",
		}
	}

	accExists, err := h.accountExists(operation.ExternalAccountID)
	if err != nil {
		return err, nil
	}

	if accExists {
		return nil, &types.OperationResult{
			Result:  ErrAccountAlreadyExists,
			Message: fmt.Sprintf("account %v already exists", operation.ExternalAccountID),
		}
	}

	res, err := h.db.CreateAccount(&models.Account{
		ExternalID: operation.ExternalAccountID,
		Balance:    0,
	})

	if err != nil || res == nil {
		return err, nil
	}

	return nil, &types.OperationResult{
		Result:  Ok,
		Message: fmt.Sprintf("Account %v created", res.ExternalID),
	}
}
