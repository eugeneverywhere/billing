package dispatcher

import (
	"encoding/json"
	"fmt"
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/db"
	"github.com/eugeneverywhere/billing/handler"
	"github.com/eugeneverywhere/billing/handler/sync"
	"github.com/eugeneverywhere/billing/rabbit"
	"github.com/eugeneverywhere/billing/types"
	"github.com/lillilli/logger"
)

type Dispatcher interface {
	Dispatch(rawMsgPayload []byte)
}

type dispatcher struct {
	log           logger.Logger
	sender        rabbit.Sender
	errSender     rabbit.Sender
	db            db.DB
	h             handler.Handler
	accountsCache cache.AccountsCache
	accountMutex  *sync.Kmutex
}

func NewDispatcher(handler handler.Handler, sender rabbit.Sender, errSender rabbit.Sender) Dispatcher {
	return &dispatcher{
		log:          logger.NewLogger("dispatcher"),
		h:            handler,
		sender:       sender,
		errSender:    errSender,
		accountMutex: sync.New(),
	}
}

func (d *dispatcher) Dispatch(rawMsgPayload []byte) {
	operationData := new(types.Operation)

	if err := json.Unmarshal(rawMsgPayload, &operationData); err != nil {
		d.log.Errorf("Can't parse operation %q: %v", string(rawMsgPayload), err)
		return
	}

	switch operationData.Code {
	case types.OpCreateAccount:
		go d.handleAccountCreate(rawMsgPayload)
	case types.OpAddAmount:
		go d.handleAddAmount(rawMsgPayload)
	case types.OpTransfer:
		go d.handleTransfer(rawMsgPayload)
	default:
		go d.sendError(&types.OperationResult{
			Operation: operationData,
			Result:    types.ErrUnknownOperationCode,
			Message:   "unknown op code",
		})
	}
}

func (d *dispatcher) handleAccountCreate(rawOperation []byte) {
	createAccountData := new(types.CreateAccount)
	if err := json.Unmarshal(rawOperation, &createAccountData); err != nil {
		d.log.Errorf("Can't parse create account operation %q: %v", string(rawOperation), err)
		go d.sendError(&types.OperationResult{
			Operation: createAccountData.Operation,
			Result:    types.ErrWrongFormat,
			Message:   fmt.Sprintf("%v", err),
		})
		return
	}

	d.log.Debugf("Handling: %v", createAccountData)

	err, result := d.h.CreateAccount(createAccountData)

	if err != nil || result == nil || result.Result != types.Ok {
		d.log.Errorf("Account creation failed for id %v: %v",
			createAccountData.ExternalAccountID, err)
		if result == nil {
			go d.sendError(&types.OperationResult{
				Operation: createAccountData.Operation,
				Result:    types.ErrInternal,
				Message:   "internal error",
			})
			return
		}
		result.Operation = createAccountData.Operation
		go d.sendError(result)
	}
}

func (d *dispatcher) handleAddAmount(rawOperation []byte) {
	addAmountData := new(types.AddAmount)
	if err := json.Unmarshal(rawOperation, &addAmountData); err != nil {
		d.log.Errorf("Can't parse add amount operation %q: %v", string(rawOperation), err)
		go d.sendError(&types.OperationResult{
			Operation: addAmountData.Operation,
			Result:    types.ErrWrongFormat,
			Message:   fmt.Sprintf("%v", err),
		})
		return
	}

	d.log.Debugf("Handling: %v", addAmountData)

	err, result := d.h.AddAmount(addAmountData)
	if err != nil || result == nil || result.Result != types.Ok {
		d.log.Errorf("Adding amount failed for id %v: %v",
			addAmountData.ExternalAccountID, err)
		if result == nil {
			go d.sendError(&types.OperationResult{
				Operation: addAmountData.Operation,
				Result:    types.ErrInternal,
				Message:   "internal error",
			})
			return
		}
		result.Operation = addAmountData.Operation
		go d.sendError(result)
	}

}

func (d *dispatcher) handleTransfer(rawOperation []byte) {
	transferData := new(types.TransferAmount)
	if err := json.Unmarshal(rawOperation, &transferData); err != nil {
		d.log.Errorf("Can't parse transfer operation %q: %v", string(rawOperation), err)
		go d.sendError(&types.OperationResult{
			Operation: transferData.Operation,
			Result:    types.ErrWrongFormat,
			Message:   fmt.Sprintf("%v", err),
		})
		return
	}
	d.log.Debugf("Handling: %v", transferData)

	err, result := d.h.TransferAmount(transferData)
	if err != nil || result == nil || result.Result != types.Ok {
		d.log.Errorf("Transfer failed for %v -> %v: %v",
			transferData.Source, transferData.Target, err)
		if result == nil {
			go d.sendError(&types.OperationResult{
				Operation: transferData.Operation,
				Result:    types.ErrInternal,
				Message:   "internal error",
			})
			return
		}
		result.Operation = transferData.Operation
		go d.sendError(result)
	}
}

func (d *dispatcher) sendResult(result *types.OperationResult) {
	err := d.sender.Send(result)
	if err != nil {
		d.log.Errorf("sending result message to rabbit failed: %v", err)
	}
}

func (d *dispatcher) sendError(result *types.OperationResult) {
	err := d.errSender.Send(result)
	if err != nil {
		d.log.Errorf("sending error message to rabbit failed: %v", err)
	}
}
