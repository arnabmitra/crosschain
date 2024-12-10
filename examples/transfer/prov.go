package main

// Cosmos specific example of sending a transaction using raw signing.

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cordialsys/crosschain/chain/cosmos/address"
	"github.com/cordialsys/crosschain/chain/cosmos/builder"
	"github.com/cordialsys/crosschain/chain/cosmos/tx"
	"github.com/cordialsys/crosschain/chain/cosmos/tx_input"
	"github.com/cordialsys/crosschain/factory"
	"github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func main() {
	ctx := context.Background()
	decodedPublicKey, errDecode := hex.DecodeString("0306cf525f0366565f082c0d8b8619c002ac989cf0562393c1a0ae30c17f1b2776")
	if errDecode != nil {
		panic(errDecode)
	}
	// initialize testnet crosschain
	xc := factory.NewNotMainnetsFactory(&factory.FactoryOptions{})

	asset, err := xc.GetAssetConfig("", "HASH")
	if err != nil {
		panic(err)
	}
	client, _ := xc.NewClient(asset)
	// cosmos builder
	builder, err := builder.NewTxBuilder(asset)

	if err != nil {
		panic("unsupported asset: " + err.Error())
	}

	from := xc.MustAddress(asset, "tp14qvkzter97k5jdhpzwpsg3u7nrendk9nvvyhqx")
	to := xc.MustAddress(asset, "tp1uywe3m7uknt8wkj78l5xar9exsthh3l3kzkuxe")
	input := tx_input.NewTxInput()
	input.AssetType = tx_input.BANK
	input.GasPrice = 19200
	input.ChainId = "pio-testnet-1"
	input.AccountNumber = 229999
	amount := xc.MustAmountBlockchain(asset, "0.001")

	xcTx, err := builder.NewTransfer(from, to, amount, input)
	cosmosTx := xcTx.(*tx.Tx).CosmosTx.(types.FeeTx)
	if err != nil {
		panic("could not create transfer object: " + err.Error())
	}
	cosmosTxConfig := builder.CosmosTxConfig
	cosmosBuilder := builder.CosmosTxBuilder
	msgs := cosmosTx.GetMsgs()
	err = cosmosBuilder.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}

	cosmosBuilder.SetGasLimit(cosmosTx.GetGas())
	cosmosBuilder.SetFeeAmount(cosmosTx.GetFee())

	fmt.Printf("the fee is %v \n", cosmosTx.GetFee())
	sigMode := signingtypes.SignMode_SIGN_MODE_DIRECT
	sigsV2 := []signingtypes.SignatureV2{
		{
			PubKey: address.GetPublicKey(asset.GetChain(), decodedPublicKey),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  sigMode,
				Signature: nil,
			},
			Sequence: input.Sequence,
		},
	}

	err = cosmosBuilder.SetSignatures(sigsV2...)
	if err != nil {
		panic(err)
	}

	chainId := input.ChainId
	if chainId == "" {
		chainId = asset.GetChain().ChainIDStr
	}

	signerData := signing.SignerData{
		AccountNumber: input.AccountNumber,
		ChainID:       chainId,
		Sequence:      input.Sequence,
	}

	sighashData, err := cosmosTxConfig.SignModeHandler().GetSignBytes(sigMode, signerData, cosmosBuilder.GetTx())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Payload to be signed: %s\n", hex.EncodeToString(sighashData))
	sighash := tx.GetSighash(asset.GetChain(), sighashData)

	txToBroadcast := &tx.Tx{
		CosmosTx:        cosmosBuilder.GetTx(),
		ParsedTransfers: msgs,
		CosmosTxBuilder: cosmosBuilder,
		CosmosTxEncoder: cosmosTxConfig.TxEncoder(),
		SigsV2:          sigsV2,
		TxDataToSign:    sighash,
	}

	// Paste in the signature of the payload here:
	signature, err := hex.DecodeString("ceb5b8d458837978cd62a086cde75dc77baf17bf0215388e91e342adfa80c2b1744ea6716fe81475a3af1f52a2860461e8f90e4365da08268b61e9b35ab2db8000")
	if err != nil {
		panic(err)
	}
	err = txToBroadcast.AddSignatures(signature)
	if err != nil {
		panic(err)
	}
	// submit the tx, wait a bit, fetch the tx info
	// (network needed)
	fmt.Printf("tx id: %s\n", txToBroadcast.Hash())
	err = client.SubmitTx(ctx, txToBroadcast)
	if err != nil {
		panic(err)
	}

}
