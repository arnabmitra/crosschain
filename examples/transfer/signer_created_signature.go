package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cordialsys/crosschain"
	"github.com/cordialsys/crosschain/factory"
)

func main() {
	// initialize crosschain
	xc := factory.NewNotMainnetsFactory(&factory.FactoryOptions{})
	ctx := context.Background()

	// get asset model, including config data
	// asset is used to create client, builder, signer, etc.
	asset, err := xc.GetAssetConfig("", "HASH")
	if err != nil {
		panic("unsupported asset: " + err.Error())
	}

	if err != nil {
		panic("could not create public key: " + err.Error())
	}

	to := xc.MustAddress(asset, "tp1uywe3m7uknt8wkj78l5xar9exsthh3l3kzkuxe")
	from := xc.MustAddress(asset, "tp1splvpmc0du6qstfk0j808jygj3p7sjwm04sq77")
	amount := xc.MustAmountBlockchain(asset, "0.001")

	// to create a tx, we typically need some input from the blockchain
	// e.g., nonce for Ethereum, recent block for Solana, gas data, ...
	// (network needed)
	client, _ := xc.NewClient(asset)

	input, err := client.FetchLegacyTxInput(ctx, from, to)
	if err != nil {
		panic(err)
	}
	if inputWithPublicKey, ok := input.(crosschain.TxInputWithPublicKey); ok {
		decodedPubLey, errDecode := hex.DecodeString("033c218532c89a97bec1cfcb1b5cd5f7d383d6f3650ccf85eb8ce77d53c687369c")
		if errDecode != nil {
			panic(errDecode)
		}
		fromPublicKeyStr := base64.StdEncoding.EncodeToString(decodedPubLey)
		errPubKey := inputWithPublicKey.SetPublicKeyFromStr(fromPublicKeyStr)
		if err != nil {
			panic(errPubKey)
		}
	}
	if inputWithAmount, ok := input.(crosschain.TxInputWithAmount); ok {
		inputWithAmount.SetAmount(amount)
	}
	fmt.Printf("%+v\n", input)

	// create tx
	// (no network, no private key needed)
	builder, _ := xc.NewTxBuilder(asset)
	tx, err := builder.NewTransfer(from, to, amount, input)
	if err != nil {
		panic(err)
	}
	sighashes, err := tx.Sighashes()
	if err != nil {
		panic(err)
	}
	sighash := sighashes[0]
	fmt.Printf("%+v\n", tx)
	fmt.Printf("signing: %x\n", hex.EncodeToString(sighash))

	// sign the tx sighash with the treasury SIGNER
	signatureRaw, _ := hex.DecodeString("63f84a3b87f91515817663ae39608f79cedb93c63cb22218ecef2cca579860eb4427460bad31eb6072c17ac927e633cea7074cc1f41d3fc626c85f71fbb378e101")
	fmt.Printf("signature: %x\n", signatureRaw)

	// complete the tx by adding its signature
	// (no network, no private key needed)
	err = tx.AddSignatures(signatureRaw)
	if err != nil {
		panic(err)
	}

	// submit the tx, wait a bit, fetch the tx info
	// (network needed)
	fmt.Printf("tx id: %s\n", tx.Hash())
	err = client.SubmitTx(ctx, tx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Zzz...")
	time.Sleep(60 * time.Second)
	info, err := client.FetchTxInfo(ctx, tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", info)
}
