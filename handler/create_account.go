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

func (h *handler) CreateAccount(operation *types.CreateAccount) (*types.OperationResult, error) {
	if ContainsSpaces(operation.ExternalAccountID) {
		return &types.OperationResult{
			Result:  types.ErrSpaces,
			Message: fmt.Sprintf("blank characters not allowed: %v", operation.ExternalAccountID),
		}, nil
	}

	if operation.ExternalAccountID == "" {
		return &types.OperationResult{
			Result:  types.ErrEmptyID,
			Message: "empty id not allowed",
		}, nil
	}

	if len(operation.ExternalAccountID) > types.MaxExternalIDLength {
		return &types.OperationResult{
			Result:  types.ErrIDTooLong,
			Message: fmt.Sprintf("length of id exceeded %v", types.MaxExternalIDLength),
		}, nil
	}

	accExists, err := h.accountExists(operation.ExternalAccountID)
	if err != nil {
		return nil, err
	}

	if accExists {
		return &types.OperationResult{
			Result:  types.ErrAccountAlreadyExists,
			Message: fmt.Sprintf("account %v already exists", operation.ExternalAccountID),
		}, nil
	}

	res, err := h.db.CreateAccount(&models.Account{
		ExternalID: operation.ExternalAccountID,
		Balance:    0,
	})

	if err != nil || res == nil {
		return nil, err
	}

	return &types.OperationResult{
		Result:  types.Ok,
		Message: fmt.Sprintf("Account %v created", res.ExternalID),
	}, nil
}
