package builder

import (
	"fmt"
	"math/big"

	xc "github.com/cordialsys/crosschain"
	"github.com/cordialsys/crosschain/chain/evm/address"
	"github.com/cordialsys/crosschain/chain/evm/tx"
	"github.com/cordialsys/crosschain/chain/evm/tx_input"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

var DefaultMaxTipCapGwei uint64 = 5

type GethTxBuilder interface {
	BuildTxWithPayload(chain *xc.ChainConfig, to xc.Address, value xc.AmountBlockchain, data []byte, input xc.TxInput) (xc.Tx, error)
}

// supports evm after london merge
type EvmTxBuilder struct {
}

var _ GethTxBuilder = &EvmTxBuilder{}

// TxBuilder for EVM
type TxBuilder struct {
	Asset         xc.ITask
	gethTxBuilder GethTxBuilder
	// Legacy bool
}

var _ xc.TxBuilder = &TxBuilder{}

func NewEvmTxBuilder() *EvmTxBuilder {
	return &EvmTxBuilder{}
}

// NewTxBuilder creates a new EVM TxBuilder
func NewTxBuilder(asset xc.ITask) (xc.TxBuilder, error) {
	return TxBuilder{
		Asset:         asset,
		gethTxBuilder: &EvmTxBuilder{},
	}, nil
}

func (txBuilder TxBuilder) WithTxBuilder(buider GethTxBuilder) xc.TxBuilder {
	txBuilder.gethTxBuilder = buider
	return txBuilder
}

// NewTxBuilder creates a new EVM TxBuilder for legacy tx
// func NewLegacyTxBuilder(asset xc.ITask) (xc.TxBuilder, error) {
// 	return TxBuilder{
// 		Asset: asset,
// 		// Legacy: true,
// 	}, nil
// }

// NewTransfer creates a new transfer for an Asset, either native or token
func (txBuilder TxBuilder) NewTransfer(from xc.Address, to xc.Address, amount xc.AmountBlockchain, input xc.TxInput) (xc.Tx, error) {
	switch asset := txBuilder.Asset.(type) {
	case *xc.TaskConfig:
		return txBuilder.NewTask(from, to, amount, input)

	case *xc.ChainConfig:
		return txBuilder.NewNativeTransfer(from, to, amount, input)

	case *xc.TokenAssetConfig:
		return txBuilder.NewTokenTransfer(from, to, amount, input)

	default:
		// TODO this should return error
		contract := asset.GetContract()
		logrus.WithFields(logrus.Fields{
			"chain":      asset.GetChain().Chain,
			"contract":   contract,
			"asset_type": fmt.Sprintf("%T", asset),
		}).Warn("new transfer for unknown asset type")
		if contract != "" {
			return txBuilder.NewTokenTransfer(from, to, amount, input)
		} else {
			return txBuilder.NewNativeTransfer(from, to, amount, input)
		}
	}
}

// NewNativeTransfer creates a new transfer for a native asset
func (txBuilder TxBuilder) NewNativeTransfer(from xc.Address, to xc.Address, amount xc.AmountBlockchain, input xc.TxInput) (xc.Tx, error) {
	return txBuilder.gethTxBuilder.BuildTxWithPayload(txBuilder.Asset.GetChain(), to, amount, []byte{}, input)
}

// NewTokenTransfer creates a new transfer for a token asset
func (txBuilder TxBuilder) NewTokenTransfer(from xc.Address, to xc.Address, amount xc.AmountBlockchain, input xc.TxInput) (xc.Tx, error) {

	zero := xc.NewAmountBlockchainFromUint64(0)
	contract := txBuilder.Asset.GetContract()
	payload, err := BuildERC20Payload(to, amount)
	if err != nil {
		return nil, err
	}
	return txBuilder.gethTxBuilder.BuildTxWithPayload(txBuilder.Asset.GetChain(), xc.Address(contract), zero, payload, input)
}

func BuildERC20Payload(to xc.Address, amount xc.AmountBlockchain) ([]byte, error) {
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	toAddress, err := address.FromHex(to)
	if err != nil {
		return nil, err
	}
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	paddedAmount := common.LeftPadBytes(amount.Int().Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data, nil
}

func (*EvmTxBuilder) BuildTxWithPayload(chain *xc.ChainConfig, to xc.Address, value xc.AmountBlockchain, data []byte, inputRaw xc.TxInput) (xc.Tx, error) {
	address, err := address.FromHex(to)
	if err != nil {
		return nil, err
	}

	input := inputRaw.(*tx_input.TxInput)
	var chainId *big.Int = input.ChainId.Int()
	if input.ChainId.Uint64() == 0 {
		chainId = new(big.Int).SetInt64(chain.ChainID)
	}

	// Protection from setting very high gas tip
	maxTipGwei := uint64(chain.ChainMaxGasPrice)
	if maxTipGwei == 0 {
		maxTipGwei = DefaultMaxTipCapGwei
	}
	maxTipWei := GweiToWei(maxTipGwei)
	gasTipCap := input.GasTipCap

	if gasTipCap.Cmp(&maxTipWei) > 0 {
		// limit to max
		gasTipCap = maxTipWei
	}

	return &tx.Tx{
		EthTx: types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     input.Nonce,
			GasTipCap: gasTipCap.Int(),
			GasFeeCap: input.GasFeeCap.Int(),
			Gas:       input.GasLimit,
			To:        &address,
			Value:     value.Int(),
			Data:      data,
		}),
		Signer: types.LatestSignerForChainID(chainId),
	}, nil
}

func GweiToWei(gwei uint64) xc.AmountBlockchain {
	bigGwei := big.NewInt(int64(gwei))

	ten := big.NewInt(10)
	nine := big.NewInt(9)
	factor := big.NewInt(0).Exp(ten, nine, nil)

	bigGwei.Mul(bigGwei, factor)
	return xc.AmountBlockchain(*bigGwei)
}