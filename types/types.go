package types

type Operation struct {
	Code int `json:"op_code"`
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
