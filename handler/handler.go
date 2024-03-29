package handler

import (
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/db"
	"github.com/eugeneverywhere/billing/handler/sync"
	"github.com/eugeneverywhere/billing/types"
	"github.com/lillilli/logger"
)

type Handler interface {
	CreateAccount(operation *types.CreateAccount) (*types.OperationResult, error)
	AddAmount(addAmount *types.AddAmount) (*types.OperationResult, error)
	TransferAmount(transfer *types.TransferAmount) (*types.OperationResult, error)
}

type handler struct {
	log           logger.Logger
	db            db.DB
	accountsCache cache.AccountsCache
	accountMutex  *sync.Kmutex
}

func NewHandler(db db.DB, accountsCache cache.AccountsCache) Handler {
	return &handler{
		log:           logger.NewLogger("handler"),
		db:            db,
		accountsCache: accountsCache,
		accountMutex:  sync.New(),
	}
}
