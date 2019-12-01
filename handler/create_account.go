package handler

import (
	"fmt"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/eugeneverywhere/billing/types"
)

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
			Result:  types.ErrSpaces,
			Message: fmt.Sprintf("blank characters not allowed: %v", operation.ExternalAccountID),
		}
	}

	if operation.ExternalAccountID == "" {
		return nil, &types.OperationResult{
			Result:  types.ErrEmptyID,
			Message: "empty id not allowed",
		}
	}

	accExists, err := h.accountExists(operation.ExternalAccountID)
	if err != nil {
		return err, nil
	}

	if accExists {
		return nil, &types.OperationResult{
			Result:  types.ErrAccountAlreadyExists,
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
		Result:  types.Ok,
		Message: fmt.Sprintf("Account %v created", res.ExternalID),
	}
}
