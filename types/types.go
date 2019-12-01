package types

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

type Operation struct {
	ConsumerID  int `json:"cons_id"`
	OperationID int `json:"op_id"`
	Code        int `json:"op_code"`
}

type OperationResult struct {
	*Operation
	Result  int    `json:"res_code"`
	Message string `json:"msg"`
}

type CreateAccount struct {
	*Operation
	ExternalAccountID string `json:"acc_id"`
}

type CreateAccountResult struct {
	*CreateAccount
	*OperationResult
}

type AddAmount struct {
	*Operation
	ExternalAccountID string  `json:"acc_id"`
	Amount            float64 `json:"amount"`
}

type AddAmountResult struct {
	*AddAmount
	*OperationResult
}

type TransferAmount struct {
	*Operation
	Amount float64 `json:"amount"`
	Source string  `json:"src_id"`
	Target string  `json:"tgt_id"`
}

type TransferAmountResult struct {
	*TransferAmount
	*OperationResult
}
