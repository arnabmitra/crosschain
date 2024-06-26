package solana

import (
	"fmt"

	xc "github.com/cordialsys/crosschain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

func (s *SolanaTestSuite) TestNewTxBuilder() {
	require := s.Require()
	builder, err := NewTxBuilder(&xc.TokenAssetConfig{Asset: "USDC"})
	require.Nil(err)
	require.NotNil(builder)
	require.Equal("USDC", builder.(TxBuilder).Asset.(*xc.TokenAssetConfig).Asset)
}

func (s *SolanaTestSuite) TestNewNativeTransfer() {
	require := s.Require()
	builder, _ := NewTxBuilder(&xc.ChainConfig{})
	from := xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	to := xc.Address("BWbmXj5ckAaWCAtzMZ97qnJhBAKegoXtgNrv9BUpAB11")
	amount := xc.NewAmountBlockchainFromUint64(1200000) // 1.2 SOL
	input := &TxInput{}
	tx, err := builder.(xc.TxTokenBuilder).NewNativeTransfer(from, to, amount, input)
	require.Nil(err)
	require.NotNil(tx)
	solTx := tx.(*Tx).SolTx
	require.Equal(0, len(solTx.Signatures))
	require.Equal(1, len(solTx.Message.Instructions))
	require.Equal(uint16(0x2), solTx.Message.Instructions[0].ProgramIDIndex) // system tx
}

func (s *SolanaTestSuite) TestNewNativeTransferErr() {
	require := s.Require()
	builder, _ := NewTxBuilder(&xc.ChainConfig{})

	from := xc.Address("from") // fails on parsing from
	to := xc.Address("to")
	amount := xc.AmountBlockchain{}
	input := &TxInput{}
	tx, err := builder.(xc.TxTokenBuilder).NewNativeTransfer(from, to, amount, input)
	require.Nil(tx)
	require.EqualError(err, "invalid length, expected 32, got 3")

	from = xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	// fails on parsing to
	tx, err = builder.(xc.TxTokenBuilder).NewNativeTransfer(from, to, amount, input)
	require.Nil(tx)
	require.EqualError(err, "invalid length, expected 32, got 2")
}

func (s *SolanaTestSuite) TestNewTokenTransfer() {
	require := s.Require()
	contract := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	builder, _ := NewTxBuilder(&xc.TokenAssetConfig{
		Contract:    contract,
		Decimals:    6,
		ChainConfig: &xc.ChainConfig{},
	})
	from := xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	to := xc.Address("BWbmXj5ckAaWCAtzMZ97qnJhBAKegoXtgNrv9BUpAB11")
	amount := xc.NewAmountBlockchainFromUint64(1200000) // 1.2 USDC

	ataToStr, _ := FindAssociatedTokenAddress(string(to), string(contract), solana.TokenProgramID)
	ataTo := solana.MustPublicKeyFromBase58(ataToStr)

	// transfer to existing ATA
	input := &TxInput{}
	tx, err := builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(err)
	require.NotNil(tx)
	solTx := tx.(*Tx).SolTx
	require.Equal(0, len(solTx.Signatures))
	require.Equal(1, len(solTx.Message.Instructions))
	require.Equal(uint16(0x4), solTx.Message.Instructions[0].ProgramIDIndex) // token tx
	require.Equal(ataTo, solTx.Message.AccountKeys[2])                       // destination

	// transfer to non-existing ATA: create
	input = &TxInput{ShouldCreateATA: true}
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(err)
	require.NotNil(tx)
	solTx = tx.(*Tx).SolTx
	require.Equal(0, len(solTx.Signatures))
	require.Equal(2, len(solTx.Message.Instructions))
	require.Equal(uint16(0x7), solTx.Message.Instructions[0].ProgramIDIndex)
	require.Equal(uint16(0x8), solTx.Message.Instructions[1].ProgramIDIndex)
	require.Equal(ataTo, solTx.Message.AccountKeys[1])

	// transfer directly to ATA
	to = xc.Address(ataToStr)
	input = &TxInput{ToIsATA: true}
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(err)
	require.NotNil(tx)
	solTx = tx.(*Tx).SolTx
	require.Equal(0, len(solTx.Signatures))
	require.Equal(1, len(solTx.Message.Instructions))
	require.Equal(uint16(0x4), solTx.Message.Instructions[0].ProgramIDIndex) // token tx
	require.Equal(ataTo, solTx.Message.AccountKeys[2])                       // destination

	// invalid: direct to ATA, but ToIsATA: false
	to = xc.Address(ataToStr)
	input = &TxInput{ToIsATA: false}
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(err)
	require.NotNil(tx)
	solTx = tx.(*Tx).SolTx
	require.Equal(0, len(solTx.Signatures))
	require.Equal(1, len(solTx.Message.Instructions))
	require.Equal(uint16(0x4), solTx.Message.Instructions[0].ProgramIDIndex) // token tx
	require.NotEqual(ataTo, solTx.Message.AccountKeys[2])                    // destination
}

func validateTransferChecked(tx *solana.Transaction, instr *solana.CompiledInstruction) (*token.TransferChecked, error) {
	accs, _ := instr.ResolveInstructionAccounts(&tx.Message)
	inst, _ := token.DecodeInstruction(accs, instr.Data)
	transferChecked := *inst.Impl.(*token.TransferChecked)
	if len(transferChecked.Signers) > 0 {
		return &transferChecked, fmt.Errorf("should not send multisig transfers")
	}
	return &transferChecked, nil
}
func getTokenTransferAmount(tx *solana.Transaction, instr *solana.CompiledInstruction) uint64 {
	transferChecked, err := validateTransferChecked(tx, instr)
	if err != nil {
		panic(err)
	}
	return *transferChecked.Amount
}

func (s *SolanaTestSuite) TestNewMultiTokenTransfer() {
	require := s.Require()
	contract := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	builder, _ := NewTxBuilder(&xc.TokenAssetConfig{
		Contract:    contract,
		Decimals:    6,
		ChainConfig: &xc.ChainConfig{},
	})
	from := xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	to := xc.Address("BWbmXj5ckAaWCAtzMZ97qnJhBAKegoXtgNrv9BUpAB11")
	amountTooBig := xc.NewAmountBlockchainFromUint64(500)
	amountExact := xc.NewAmountBlockchainFromUint64(300)
	amountSmall1 := xc.NewAmountBlockchainFromUint64(100)
	amountSmall2 := xc.NewAmountBlockchainFromUint64(150)
	amountSmall3 := xc.NewAmountBlockchainFromUint64(200)

	ataToStr, err := FindAssociatedTokenAddress(string(to), string(contract), solana.TokenProgramID)
	require.NoError(err)
	ataTo := solana.MustPublicKeyFromBase58(ataToStr)

	// transfer to existing ATA
	input := &TxInput{
		SourceTokenAccounts: []*TokenAccount{
			{
				Account: solana.PublicKey{},
				Balance: xc.NewAmountBlockchainFromUint64(100),
			},
			{
				Account: solana.PublicKey{},
				Balance: xc.NewAmountBlockchainFromUint64(100),
			},
			{
				Account: solana.PublicKey{},
				Balance: xc.NewAmountBlockchainFromUint64(100),
			},
		},
	}
	_, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amountTooBig, input)
	require.ErrorContains(err, "cannot send")

	tx, err := builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amountExact, input)
	require.NoError(err)
	solTx := tx.(*Tx).SolTx

	_, err = validateTransferChecked(solTx, &solTx.Message.Instructions[0])
	require.NoError(err)

	require.Equal(uint16(0x4), solTx.Message.Instructions[0].ProgramIDIndex) // token tx
	require.Equal(ataTo, solTx.Message.AccountKeys[2])                       // destination
	require.Equal(3, len(solTx.Message.Instructions))
	// exactAmount should have 3 instructions, 100 amount each
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[0]))
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[1]))
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[2]))

	// amountSmall1 should just have 1 instruction (fits 1 token balance exact)
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amountSmall1, input)
	require.NoError(err)
	solTx = tx.(*Tx).SolTx
	require.Equal(1, len(solTx.Message.Instructions))
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[0]))

	// amountSmall2 should just have 2 instruction (first 100, second 50)
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amountSmall2, input)
	require.NoError(err)
	solTx = tx.(*Tx).SolTx
	require.Equal(2, len(solTx.Message.Instructions))
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[0]))
	require.EqualValues(50, getTokenTransferAmount(solTx, &solTx.Message.Instructions[1]))

	// amountSmall3 should just have 3 instruction (first 100, second 100)
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amountSmall3, input)
	require.NoError(err)
	solTx = tx.(*Tx).SolTx
	require.Equal(2, len(solTx.Message.Instructions))
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[0]))
	require.EqualValues(100, getTokenTransferAmount(solTx, &solTx.Message.Instructions[1]))

}

func (s *SolanaTestSuite) TestNewTokenTransferErr() {
	require := s.Require()

	// invalid asset
	builder, _ := NewTxBuilder(&xc.ChainConfig{})
	from := xc.Address("from")
	to := xc.Address("to")
	amount := xc.AmountBlockchain{}
	input := &TxInput{}
	tx, err := builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(tx)
	require.EqualError(err, "asset does not have a contract")

	// invalid from, to
	builder, _ = NewTxBuilder(&xc.TokenAssetConfig{
		Contract: "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
		Decimals: 6,
	})
	from = xc.Address("from")
	to = xc.Address("to")
	amount = xc.AmountBlockchain{}
	input = &TxInput{}
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(tx)
	require.EqualError(err, "invalid length, expected 32, got 3")

	from = xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(tx)
	require.EqualError(err, "invalid length, expected 32, got 2")

	// invalid asset config
	builder, _ = NewTxBuilder(&xc.TokenAssetConfig{
		Contract: "contract",
		Decimals: 6,
	})
	tx, err = builder.(xc.TxTokenBuilder).NewTokenTransfer(from, to, amount, input)
	require.Nil(tx)
	require.EqualError(err, "invalid length, expected 32, got 6")
}

func (s *SolanaTestSuite) TestNewTransfer() {
	require := s.Require()
	builder, _ := NewTxBuilder(&xc.ChainConfig{})
	from := xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	to := xc.Address("BWbmXj5ckAaWCAtzMZ97qnJhBAKegoXtgNrv9BUpAB11")
	amount := xc.NewAmountBlockchainFromUint64(1200000) // 1.2 SOL
	input := &TxInput{}
	tx, err := builder.NewTransfer(from, to, amount, input)
	require.Nil(err)
	require.NotNil(tx)
	solTx := tx.(*Tx).SolTx
	require.Equal(0, len(solTx.Signatures))
	require.Equal(1, len(solTx.Message.Instructions))
	require.Equal(uint16(0x2), solTx.Message.Instructions[0].ProgramIDIndex) // system tx
}

func (s *SolanaTestSuite) TestNewTransferAsToken() {
	require := s.Require()
	builder, _ := NewTxBuilder(&xc.TokenAssetConfig{
		Contract:    "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
		Decimals:    6,
		ChainConfig: &xc.ChainConfig{},
	})
	from := xc.Address("Hzn3n914JaSpnxo5mBbmuCDmGL6mxWN9Ac2HzEXFSGtb")
	to := xc.Address("BWbmXj5ckAaWCAtzMZ97qnJhBAKegoXtgNrv9BUpAB11")
	amount := xc.NewAmountBlockchainFromUint64(1200000) // 1.2 SOL

	type testcase struct {
		txInput               *TxInput
		expectedSourceAccount string
	}
	testcases := []testcase{
		{
			txInput: &TxInput{
				RecentBlockHash: solana.HashFromBytes([]byte{1, 2, 3, 4}),
			},
			expectedSourceAccount: "DvSgNMRxVSMBpLp4hZeBrmQo8ZRFne72actTZ3PYE3AA",
		},
		{
			txInput: &TxInput{
				RecentBlockHash: solana.HashFromBytes([]byte{1, 2, 3, 4}),
				SourceTokenAccounts: []*TokenAccount{
					{
						Account: solana.MustPublicKeyFromBase58("gCr8Xc43gEKntp7pjsBNq8qFHeUUdie2D7TrfbzPMJP"),
					},
				},
			},
			// should use new source account specified in txInput
			expectedSourceAccount: "gCr8Xc43gEKntp7pjsBNq8qFHeUUdie2D7TrfbzPMJP",
		},
	}
	for _, v := range testcases {
		tx, err := builder.NewTransfer(from, to, amount, v.txInput)
		require.Nil(err)
		require.NotNil(tx)
		solTx := tx.(*Tx).SolTx
		require.Equal(0, len(solTx.Signatures))
		require.Equal(1, len(solTx.Message.Instructions))
		require.Equal(uint16(0x4), solTx.Message.Instructions[0].ProgramIDIndex) // token tx
		tokenTf, err := validateTransferChecked(solTx, &solTx.Message.Instructions[0])
		require.NoError(err)
		require.Equal(v.expectedSourceAccount, tokenTf.Accounts[0].PublicKey.String())
	}
}
