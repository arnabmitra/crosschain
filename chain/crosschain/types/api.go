package types

import (
	xc "github.com/jumpcrypto/crosschain"
)

type ApiResult interface{}

type ApiResponse struct {
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	// wrap result inside a property of the response, so we can sign it
	Result ApiResult `json:"result,omitempty"`
}

type ChainReq struct {
	Chain string `json:"chain"`
}

type AssetReq struct {
	*ChainReq
	Asset    string `json:"asset,omitempty"`
	Contract string `json:"contract,omitempty"`
	Decimals string `json:"decimals,omitempty"`
}

type BalanceReq struct {
	*AssetReq
	Address string `json:"address"`
}

type BalanceRes struct {
	Object string `json:"object"`
	*BalanceReq
	Balance    xc.AmountHumanReadable `json:"balance"`
	BalanceRaw xc.AmountBlockchain    `json:"balance_raw"`
}

type TxInputReq struct {
	*AssetReq
	From string `json:"from"`
	To   string `json:"to"`
}

type TxInputRes struct {
	Object string `json:"object"`
	*TxInputReq
	xc.TxInput `json:"raw_tx_input,omitempty"`
}

type TxInfoReq struct {
	*AssetReq
	TxHash string `json:"tx_hash"`
}

type TxInfoRes struct {
	Object string `json:"object"`
	*TxInfoReq
	xc.TxInfo `json:"tx_info,omitempty"`
}

type SubmitTxReq struct {
	*ChainReq
	TxData []byte `json:"tx_data"`
}

type SubmitTxRes struct {
	Object string `json:"object"`
	*SubmitTxReq
}

var _ xc.Tx = &SubmitTxReq{}

func (tx *SubmitTxReq) Hash() xc.TxHash {
	panic("not implemented")
}
func (tx *SubmitTxReq) Sighashes() ([]xc.TxDataToSign, error) {
	panic("not implemented")
}
func (tx *SubmitTxReq) AddSignatures(...xc.TxSignature) error {
	panic("not implemented")
}
func (tx *SubmitTxReq) Serialize() ([]byte, error) {
	return tx.TxData, nil
}