package evm_legacy

import (
	"fmt"
	"math/big"

	xc "github.com/cordialsys/crosschain"
	"github.com/cordialsys/crosschain/chain/evm"
	"github.com/ethereum/go-ethereum/core/types"
)

var DefaultMaxTipCapGwei uint64 = 5

// TxBuilder for EVM
type TxBuilder = evm.TxBuilder

var _ xc.TxBuilder = &TxBuilder{}

// NewTxBuilder creates a new EVM TxBuilder
func NewTxBuilder(asset xc.ITask) (xc.TxBuilder, error) {
	builder, err := evm.NewTxBuilder(asset)
	if err != nil {
		return builder, err
	}
	fmt.Println("-- new legacy tx builder:: ", asset.GetChain().Chain)
	return builder.(evm.TxBuilder).WithTxBuilder(&LegacyEvmTxBuilder{}), nil
}

// supports evm before london merge
type LegacyEvmTxBuilder struct {
}

var _ evm.GethTxBuilder = &LegacyEvmTxBuilder{}

func (*LegacyEvmTxBuilder) BuildTxWithPayload(chain *xc.ChainConfig, to xc.Address, value xc.AmountBlockchain, data []byte, inputRaw xc.TxInput) (xc.Tx, error) {
	address, err := evm.HexToAddress(to)
	if err != nil {
		return nil, err
	}
	chainID := new(big.Int).SetInt64(chain.ChainID)
	input := inputRaw.(*TxInput)
	// Protection from setting very high gas tip
	// TODO

	return &Tx{
		EthTx: types.NewTransaction(
			input.Nonce,
			address,
			value.Int(),
			input.GasLimit,
			input.GasPrice.Int(),
			data,
		),
		Signer: types.LatestSignerForChainID(chainID),
	}, nil
}
