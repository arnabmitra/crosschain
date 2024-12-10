package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cordialsys/crosschain"
	"github.com/cordialsys/crosschain/factory"
	"time"
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

	// set your own private key and address
	// you can get them, for example, from your Phantom wallet
	//privateKeyInput := os.Getenv("PRIVATE_KEY")
	//if privateKeyInput == "" {
	//	log.Fatalln("must set env PRIVATE_KEY")
	//}

	//signer, _ := xc.NewSigner(asset, privateKeyInput)
	publicKey, err := hex.DecodeString("033c218532c89a97bec1cfcb1b5cd5f7d383d6f3650ccf85eb8ce77d53c687369c")
	if err != nil {
		panic("could not create public key: " + err.Error())
	}

	addressBuilder, err := xc.NewAddressBuilder(asset)
	if err != nil {
		panic("could not create address builder: " + err.Error())
	}

	from, err := addressBuilder.GetAddressFromPublicKey(publicKey)
	if err != nil {
		panic("could create from address: " + err.Error())
	}
	fmt.Println("from:", from)
	to := xc.MustAddress(asset, "tp1uywe3m7uknt8wkj78l5xar9exsthh3l3kzkuxe")
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
		fromPublicKeyStr := base64.StdEncoding.EncodeToString(publicKey)
		inputWithPublicKey.SetPublicKeyFromStr(fromPublicKeyStr)
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

	// sign the tx sighash
	//signature, err := signer.Sign(sighash)
	//if err != nil {
	//	panic(err)
	//}
	signature, err := hex.DecodeString("68426dc7f8c01bbf41182970cd951e2a5a4e5a51bf0b733a0802ff2735aa5c55508ab7ce895e6f4221255ba858628b138b483420e24318bff38fddd6164313a001")
	//signature := []byte("68426dc7f8c01bbf41182970cd951e2a5a4e5a51bf0b733a0802ff2735aa5c55508ab7ce895e6f4221255ba858628b138b483420e24318bff38fddd6164313a001")
	fmt.Printf("signature: %x\n", signature)

	// complete the tx by adding its signature
	// (no network, no private key needed)
	err = tx.AddSignatures(signature)
	if err != nil {
		panic(err)
	}

	//submit the tx, wait a bit, fetch the tx info
	//(network needed)
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
