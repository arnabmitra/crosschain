package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jumpcrypto/crosschain"
	"github.com/jumpcrypto/crosschain/factory"
)

func main() {
	// initialize crosschain
	xc := factory.NewDefaultFactory()
	ctx := context.Background()

	// get asset model, including config data
	// asset is used to create client, builder, signer, etc.
	asset, err := xc.GetAssetConfig("", "CHZ")
	if err != nil {
		panic("unsupported asset: " + err.Error())
	}

	// set your own private key and address
	// you can get them, for example, from your Phantom wallet
	privateKeyInput := os.Getenv("PRIVATE_KEY")
	if privateKeyInput == "" {
		log.Fatalln("must set env PRIVATE_KEY")
	}

	fromPrivateKey := xc.MustPrivateKey(asset, privateKeyInput)

	signer, _ := xc.NewSigner(asset)
	publicKey, err := signer.PublicKey(fromPrivateKey)
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
	to := xc.MustAddress(asset, "inj106cr7vek048k72g960l0rqpsy273ketzcx59f2")
	amount := xc.MustAmountBlockchain(asset, "0.001")

	// to create a tx, we typically need some input from the blockchain
	// e.g., nonce for Ethereum, recent block for Solana, gas data, ...
	// (network needed)
	client, _ := xc.NewClient(asset)

	input, err := client.FetchTxInput(ctx, from, to)
	if err != nil {
		panic(err)
	}
	if inputWithPublicKey, ok := input.(crosschain.TxInputWithPublicKey); ok {
		fromPublicKeyStr := base64.StdEncoding.EncodeToString(publicKey)
		inputWithPublicKey.SetPublicKeyFromStr(fromPublicKeyStr)
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
	fmt.Printf("signing: %x\n", sighash)

	// sign the tx sighash
	signature, err := signer.Sign(fromPrivateKey, sighash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %x\n", signature)

	// complete the tx by adding its signature
	// (no network, no private key needed)
	err = tx.AddSignatures(signature)
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
