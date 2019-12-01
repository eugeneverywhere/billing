package handler

import (
	"encoding/json"
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/db"
	"github.com/eugeneverywhere/billing/handler/sync"
	"github.com/eugeneverywhere/billing/rabbit"
	"github.com/eugeneverywhere/billing/types"
	"github.com/lillilli/logger"
)

const (
	OpCreateAccount = 1
	OpAddAmount     = 2
	OpTransfer      = 3

	Ok = 0

	ErrWrongFormat          = 101
	ErrUnknownOperationCode = 102
	ErrAccountAlreadyExists = 103
	ErrAccountDoesNotExist  = 104
	ErrInsufficient         = 105
	ErrEmptyID              = 106
	ErrSpaces               = 107
	ErrNonPositive          = 108

	ErrInternal = -1
)

type Handler interface {
	Dispatch(rawMsgPayload []byte)
}

type handler struct {
	log           logger.Logger
	db            db.DB
	accountsCache cache.AccountsCache
	sender        rabbit.Sender
	errSender     rabbit.Sender
	accountMutex  *sync.Kmutex
}

func NewHandler(db db.DB, accountsCache cache.AccountsCache, sender rabbit.Sender, errSender rabbit.Sender) Handler {
	return &handler{
		log:           logger.NewLogger("handler"),
		db:            db,
		accountsCache: accountsCache,
		sender:        sender,
		errSender:     errSender,
		accountMutex:  sync.New(),
	}
}

func (h *handler) Dispatch(rawMsgPayload []byte) {
	operationData := new(types.Operation)

	if err := json.Unmarshal(rawMsgPayload, &operationData); err != nil {
		h.log.Errorf("Can't parse operation %q: %v", string(rawMsgPayload), err)
		return
	}

	switch operationData.Code {
	case OpCreateAccount:
		go h.handleAccountCreate(rawMsgPayload)
	case OpAddAmount:
		go h.handleAddAmount(rawMsgPayload)
	case OpTransfer:
		go h.handleTransfer(rawMsgPayload)
	default:
		go h.sendError(&types.OperationResult{
			Operation: &types.Operation{Code: operationData.Code},
			Result:    ErrUnknownOperationCode,
			Message:   "unknown op code",
		})
	}
}

func (h *handler) sendResult(result *types.OperationResult) {
	err := h.sender.Send(result)
	if err != nil {
		h.log.Errorf("sending result message to rabbit failed: %v", err)
	}
}

func (h *handler) sendError(result *types.OperationResult) {
	err := h.errSender.Send(result)
	if err != nil {
		h.log.Errorf("sending error message to rabbit failed: %v", err)
	}
}
