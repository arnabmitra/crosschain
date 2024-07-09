package blockbook

import (
	xc "github.com/cordialsys/crosschain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
type BlockbookStats struct {
	Coin            string `json:"coin"`
	Host            string `json:"host"`
	Version         string `json:"version"`
	GitCommit       string `json:"gitCommit"`
	BuildTime       string `json:"buildTime"`
	SyncMode        bool   `json:"syncMode"`
	InitialSync     bool   `json:"initialSync"`
	InSync          bool   `json:"inSync"`
	BestHeight      int64  `json:"bestHeight"`
	LastBlockTime   string `json:"lastBlockTime"`
	InSyncMempool   bool   `json:"inSyncMempool"`
	LastMempoolTime string `json:"lastMempoolTime"`
	MempoolSize     int    `json:"mempoolSize"`
	Decimals        int    `json:"decimals"`
	DBSize          int64  `json:"dbSize"`
	About           string `json:"about"`
}

type BackendStats struct {
	Chain           string `json:"chain"`
	Blocks          int    `json:"blocks"`
	Headers         int    `json:"headers"`
	BestBlockHash   string `json:"bestBlockHash"`
	Difficulty      string `json:"difficulty"`
	SizeOnDisk      int64  `json:"sizeOnDisk"`
	Version         string `json:"version"`
	Subversion      string `json:"subversion"`
	ProtocolVersion string `json:"protocolVersion"`
	TimeOffset      int64  `json:"timeOffset"`
	Warnings        string `json:"warnings"`
}

type StatsResponse struct {
	Blockbook BlockbookStats `json:"blockbook"`
	Backend   BackendStats   `json:"backend"`
}

type UtxoResponse []Utxo
type Utxo struct {
	TxID          string `json:"txid"`
	Vout          int    `json:"vout"`
	Value         string `json:"value"`
	Confirmations uint64 `json:"confirmations"`
	LockTime      int64  `json:"lockTime"`
	Height        int64  `json:"height"`
}

func (u Utxo) GetValue() uint64 {
	return xc.NewAmountBlockchainFromStr(u.Value).Uint64()
}
func (u Utxo) GetBlock() uint64 {
	if u.Height < 0 {
		return 0
	}
	return uint64(u.Height)
}
func (u Utxo) GetTxHash() string {
	return u.TxID
}
func (u Utxo) GetIndex() uint32 {
	return uint32(u.Vout)
}

type Vin struct {
	TxID      string   `json:"txid"`
	Vout      int      `json:"vout"`
	Sequence  uint32   `json:"sequence"`
	N         int      `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
	Value     string   `json:"value"`
	Hex       string   `json:"hex"`
}

type Vout struct {
	Value     string   `json:"value"`
	N         int      `json:"n"`
	Hex       string   `json:"hex"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
}

type TransactionResponse struct {
	TxID          string `json:"txid"`
	Version       int    `json:"version"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockHash"`
	BlockHeight   int    `json:"blockHeight"`
	Confirmations int    `json:"confirmations"`
	BlockTime     int64  `json:"blockTime"`
	Size          int    `json:"size"`
	Vsize         int    `json:"vsize"`
	Value         string `json:"value"`
	ValueIn       string `json:"valueIn"`
	Fees          string `json:"fees"`
	Hex           string `json:"hex"`
}

type EstimateFeeResponse struct {
	// This is a decimal string.  It is BTC/kilobyte.
	Result string `json:"result"`
}