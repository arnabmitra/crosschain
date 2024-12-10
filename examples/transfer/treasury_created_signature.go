package main

import (
	"context"
	"encoding/base64"
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
	decodedPublicKey, errDecode := hex.DecodeString("033c218532c89a97bec1cfcb1b5cd5f7d383d6f3650ccf85eb8ce77d53c687369c")
	if errDecode != nil {
		panic(errDecode)
	}
	// initialize crosschain
	xc := factory.NewNotMainnetsFactory(&factory.FactoryOptions{})

	asset, err := xc.GetAssetConfig("", "HASH")
	client, _ := xc.NewClient(asset)
	// cosmos builder
	builder, err := builder.NewTxBuilder(asset)

	if err != nil {
		panic("unsupported asset: " + err.Error())
	}

	to := xc.MustAddress(asset, "tp1uywe3m7uknt8wkj78l5xar9exsthh3l3kzkuxe")
	from := xc.MustAddress(asset, "tp1splvpmc0du6qstfk0j808jygj3p7sjwm04sq77")
	input := tx_input.NewTxInput()
	input.AssetType = tx_input.BANK
	input.GasPrice = 19200
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
	if err != nil {
		panic(err)
	}
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
	encoded := base64.StdEncoding.EncodeToString(sighashData)
	fmt.Printf("This is the encoded string %v \n ", encoded)
	sighash := tx.GetSighash(asset.GetChain(), sighashData)

	// THIS iS WHAT GETS SIGNED (have also tried with sigHash)
	fmt.Printf("This is what get turned into hex %s \n ", sighashData)
	encodedHex := hex.EncodeToString(sighashData)
	fmt.Println(encodedHex)

	txToBroadcast := &tx.Tx{
		CosmosTx:        cosmosBuilder.GetTx(),
		ParsedTransfers: msgs,
		CosmosTxBuilder: cosmosBuilder,
		CosmosTxEncoder: cosmosTxConfig.TxEncoder(),
		SigsV2:          sigsV2,
		TxDataToSign:    sighash,
	}

	//get this from the command line for now
	// treasury  signatures create  --key 1555 --message 0a8b010a88010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e6412680a297470313471766b7a74657239376b356a6468707a777073673375376e72656e646b396e767679687178122974703175797765336d37756b6e7438776b6a37386c3578617239657873746868336c336b7a6b7578651a100a056e68617368120731303030303030126b0a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a210306cf525f0366565f082c0d8b8619c002ac989cf0562393c1a0ae30c17f1b277612040a02080112190a130a056e68617368120a373638303030303030301080b5181a0d70696f2d746573746e65742d3120ef840e  --sign-with mpc-key --format recovery
	//{
	//	"signatures": [
	//{
	//"name": "signatures/1536",
	//"state": "signed",
	//"create_time": "2024-11-03T17:01:14Z",
	//"update_time": "2024-11-03T17:01:20Z",
	//"version": 4,
	//"format": "raw",
	//"signature": "093baf176f0660546ed61b50d0beb2d863e0acc7bc8fde9af557ee4b7dc01154653af7b7c07ce064ea35c92f0e4bfb8d1a0da83279f69b4c503ebcccac3a51ab",
	//"shares": {
	//"signers/1": "030000000300000000000000010000000000000031010000000000000033010000000000000034210000000000000003093baf176f0660546ed61b50d0beb2d863e0acc7bc8fde9af557ee4b7dc011540300000000000000010000000000000031210000000000000002b41089dd360537efc7d7940695e9c293f2480030381091e2c4eaa0474b137ffd010000000000000033210000000000000003ddedc51822db2d85c6b8f5ad23849f0ccedac5b58f223be431c9c9155a238e0c0100000000000000342100000000000000030f9acb3be90d19118bda1180fdc434b299516dd00c9fa117c480eea8356117a00096a2811b85b8aeb360c8fb6451fe5517e647c76cee920898884c74ce760f25",
	//"signers/3": "030000000300000000000000010000000000000031010000000000000033010000000000000034210000000000000003093baf176f0660546ed61b50d0beb2d863e0acc7bc8fde9af557ee4b7dc011540300000000000000010000000000000031210000000000000002b41089dd360537efc7d7940695e9c293f2480030381091e2c4eaa0474b137ffd010000000000000033210000000000000003ddedc51822db2d85c6b8f5ad23849f0ccedac5b58f223be431c9c9155a238e0c0100000000000000342100000000000000030f9acb3be90d19118bda1180fdc434b299516dd00c9fa117c480eea8356117a05955d476897bd8eac7a1af43483dd505f767c9bf4df551030df741bfbafa484a",
	//"signers/4": "030000000300000000000000010000000000000031010000000000000033010000000000000034210000000000000003093baf176f0660546ed61b50d0beb2d863e0acc7bc8fde9af557ee4b7dc011540300000000000000010000000000000031210000000000000002b41089dd360537efc7d7940695e9c293f2480030381091e2c4eaa0474b137ffd010000000000000033210000000000000003ddedc51822db2d85c6b8f5ad23849f0ccedac5b58f223be431c9c9155a238e0c0100000000000000342100000000000000030f9acb3be90d19118bda1180fdc434b299516dd00c9fa117c480eea8356117a04c436c331b6f60133e4c0360b98bb1d4a4f55bbd4816e2a89a9f2dc92cce209f"
	//},
	//"message": "Q29zQkNvZ0JDaHd2WTI5emJXOXpMbUpoYm1zdWRqRmlaWFJoTVM1TmMyZFRaVzVrRW1nS0tYUndNWE53Ykhad2JXTXdaSFUyY1hOMFptc3dhamd3T0dwNVoyb3pjRGR6YW5kdE1EUnpjVGMzRWlsMGNERjFlWGRsTTIwM2RXdHVkRGgzYTJvM09HdzFlR0Z5T1dWNGMzUm9hRE5zTTJ0NmEzVjRaUm9RQ2dWdWFHRnphQklITVRBd01EQXdNQkpyQ2s0S1Jnb2ZMMk52YzIxdmN5NWpjbmx3ZEc4dWMyVmpjREkxTm1zeExsQjFZa3RsZVJJakNpRURQQ0dGTXNpYWw3N0J6OHNiWE5YMzA0UFc4MlVNejRYcmpPZDlVOGFITnB3U0JBb0NDQUVTR1FvVENnVnVhR0Z6YUJJS056WTRNREF3TURBd01CQ0F0UmdhRFhCcGJ5MTBaWE4wYm1WMExURT0=",
	//"key": "keys/1514",
	//"digest": null,
	//"triples": [
	//316,
	//317
	//]
	//}
	//],
	//"page_size": 1,
	//"total_size": 1
	//}

	signature, err := hex.DecodeString("01bb7fbf38b83f245ebdbd743782cad5be00208bb829f6b1edd27d3d1fe8684b19ccc1cc4ccbf38c525d688205c89a209deea9d721056ae2d7c71f55e78e103c00")
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
