package newchain

import (
	"context"
	"errors"

	xc "github.com/cordialsys/crosschain"
	xclient "github.com/cordialsys/crosschain/client"
	"github.com/cordialsys/crosschain/utils"
)

// Client for Template
type Client struct {
}

var _ xclient.FullClient = &Client{}

// TxInput for Template
type TxInput struct {
	xc.TxInputEnvelope
	utils.TxPriceInput
}

func NewTxInput() *TxInput {
	return &TxInput{
		TxInputEnvelope: xc.TxInputEnvelope{
			Type: "INPUT_DRIVER_HERE",
		},
	}
}

// NewClient returns a new Template Client
func NewClient(cfgI xc.ITask) (*Client, error) {
	return &Client{}, errors.New("not implemented")
}

// FetchTxInput returns tx input for a Template tx
func (client *Client) FetchTxInput(ctx context.Context, from xc.Address, to xc.Address) (xc.TxInput, error) {
	return &TxInput{}, errors.New("not implemented")
}

// SubmitTx submits a Template tx
func (client *Client) SubmitTx(ctx context.Context, txInput xc.Tx) error {
	return errors.New("not implemented")
}

// Returns transaction info - legacy/old endpoint
func (client *Client) FetchLegacyTxInfo(ctx context.Context, txHash xc.TxHash) (xc.LegacyTxInfo, error) {
	return xc.LegacyTxInfo{}, errors.New("not implemented")
}

// Returns transaction info - new endpoint
func (client *Client) FetchTxInfo(ctx context.Context, txHash xc.TxHash) (xclient.TxInfo, error) {
	return xclient.TxInfo{}, errors.New("not implemented")
}

func (client *Client) FetchNativeBalance(ctx context.Context, address xc.Address) (xc.AmountBlockchain, error) {
	return xc.AmountBlockchain{}, errors.New("not implemented")
}

func (client *Client) FetchBalance(ctx context.Context, address xc.Address) (xc.AmountBlockchain, error) {
	return xc.AmountBlockchain{}, errors.New("not implemented")
}
