package chains

import xc "github.com/jumpcrypto/crosschain"

func init() {
	for _, chain := range Testnet {
		if chain.Net == "" {
			chain.Net = "testnet"
		}
	}
}

var Testnet = []*xc.NativeAssetConfig{

	{
		Asset:       string(xc.APTOS),
		Driver:      string(xc.DriverAptos),
		Net:         "devnet",
		URL:         "https://fullnode.devnet.aptoslabs.com",
		ChainName:   "Aptos (Devnet)",
		ExplorerURL: "https://explorer.devnet.aptos.dev",
		Decimals:    8,
	},
	{
		Asset:           string(xc.ATOM),
		Driver:          string(xc.DriverCosmos),
		URL:             "https://rpc.sentry-01.theta-testnet.polypore.xyz",
		ChainName:       "Cosmos (Theta Testnet)",
		ExplorerURL:     "https://explorer.theta-testnet.polypore.xyz",
		Decimals:        6,
		ChainIDStr:      "theta-testnet-001",
		ChainPrefix:     "cosmos",
		ChainCoin:       "uatom",
		ChainCoinHDPath: 118,
	},
	{
		Asset:       string(xc.AVAX),
		Driver:      string(xc.DriverEVM),
		URL:         "https://api.avax-test.network/ext/bc/C/rpc",
		ChainName:   "Avalanche (Fuji Testnet)",
		ExplorerURL: "https://testnet.snowtrace.io",
		IndexerUrl:  "https://api.covalenthq.com/v1",
		IndexerType: "covalent",
		Decimals:    18,
		ChainID:     43113,
	},
	{
		Asset:       string(xc.BCH),
		Driver:      string(xc.DriverBitcoin),
		URL:         "",
		ChainName:   "Bitcoin Cash (Testnet)",
		ExplorerURL: "",
		IndexerType: "none",
		Decimals:    8,
	},
	{
		Asset:       string(xc.BNB),
		Driver:      string(xc.DriverEVMLegacy),
		URL:         "https://data-seed-prebsc-1-s1.binance.org:8545",
		ChainName:   "Binance Smart Chain (Testnet)",
		ExplorerURL: "https://testnet.bscscan.com",
		Decimals:    18,
		ChainID:     97,
	},
	{
		Asset:       string(xc.BTC),
		Driver:      string(xc.DriverBitcoin),
		URL:         "https://api.blockchair.com/bitcoin/testnet",
		ChainName:   "Bitcoin (Testnet)",
		ExplorerURL: "https://blockchair.com/bitcoin/testnet",
		Auth:        "env:BLOCKCHAIR_API_TOKEN",
		Decimals:    8,
		Provider:    "blockchair",
	},
	{
		Asset:       string(xc.CHZ2),
		Driver:      string(xc.DriverEVMLegacy),
		URL:         "https://spicy-rpc.chiliz.com",
		ChainName:   "Chiliz 2.0 (testnet)",
		ExplorerURL: "https://spicy-explorer.chiliz.com",
		IndexerUrl:  "https://spicy-explorer.chiliz.com",
		IndexerType: "blockscout",
		Decimals:    18,
		ChainID:     88882,
	},
	{
		Asset:       string(xc.DOGE),
		Driver:      string(xc.DriverBitcoin),
		URL:         "",
		ChainName:   "Dogecoin (Testnet)",
		ExplorerURL: "",
		IndexerType: "none",
		Decimals:    8,
	},
	{
		Asset:       string(xc.ETC),
		Driver:      string(xc.DriverEVMLegacy),
		URL:         "https://www.ethercluster.com/mordor",
		ChainName:   "Ethereum Classic (Mordor)",
		ExplorerURL: "",
		IndexerType: "none",
		Decimals:    18,
		ChainID:     63,
	},
	{
		Asset:       string(xc.ETH),
		Driver:      string(xc.DriverEVM),
		URL:         "https://goerli.infura.io/v3",
		ChainName:   "Ethereum (Goerli)",
		ExplorerURL: "https://goerli.etherscan.io",
		Auth:        "env:INFURA_API_TOKEN",
		IndexerType: "rpc",
		Decimals:    18,
		ChainID:     5,
		Provider:    "infura",
	},
	{
		Asset:       string(xc.CELO),
		Driver:      string(xc.DriverEVM),
		URL:         "https://alfajores-forno.celo-testnet.org",
		ChainName:   "Celo (Testnet)",
		ExplorerURL: "https://alfajores-blockscout.celo-testnet.org/",
		IndexerType: "rpc",
		Decimals:    18,
		ChainID:     44787,
	},
	{
		Asset:       string(xc.FTM),
		Driver:      string(xc.DriverEVMLegacy),
		URL:         "https://rpc.testnet.fantom.network",
		ChainName:   "Fantom (Testnet)",
		ExplorerURL: "https://testnet.ftmscan.com",
		IndexerType: "none",
		Decimals:    18,
		ChainID:     4002,
	},
	{
		Asset:                string(xc.INJ),
		Driver:               string(xc.DriverCosmos),
		URL:                  "https://k8s.testnet.tm.injective.network",
		ChainName:            "Injective (Testnet)",
		ExplorerURL:          "https://testnet.explorer.injective.network",
		Decimals:             18,
		ChainGasPriceDefault: 500000000.00,
		ChainIDStr:           "injective-888",
		ChainPrefix:          "inj",
		ChainCoin:            "inj",
		ChainCoinHDPath:      60,
	},
	{
		Asset:              string(xc.LUNC),
		Driver:             string(xc.DriverCosmos),
		URL:                "",
		ChainName:          "Terra Classic (Testnet)",
		Decimals:           6,
		ChainGasMultiplier: 12.00,
		ChainIDStr:         "bombay-12",
		ChainPrefix:        "terra",
		ChainCoin:          "uluna",
		ChainCoinHDPath:    330,
	},
	{
		Asset:              string(xc.LUNA),
		Driver:             string(xc.DriverCosmos),
		URL:                "https://terra-testnet-rpc.polkachu.com",
		ChainName:          "Terra (Testnet)",
		ExplorerURL:        "https://finder.terra.money/testnet",
		Decimals:           6,
		ChainGasMultiplier: 12.00,
		ChainIDStr:         "pisco-1",
		ChainPrefix:        "terra",
		ChainCoin:          "uluna",
		ChainCoinHDPath:    330,
	},
	{
		Asset:          string(xc.MATIC),
		Driver:         string(xc.DriverEVM),
		URL:            "https://rpc-mumbai.matic.today",
		ChainName:      "Polygon (Mumbai)",
		ExplorerURL:    "https://mumbai.polygonscan.com",
		IndexerType:    "rpc",
		Decimals:       18,
		ChainID:        80001,
		ChainMaxGasTip: 120,
	},
	{
		Asset:       string(xc.ROSE),
		Driver:      string(xc.DriverEVMLegacy),
		URL:         "https://testnet.emerald.oasis.dev",
		ChainName:   "Oasis Emerald (Testnet)",
		ExplorerURL: "",
		IndexerType: "none",
		Decimals:    18,
		ChainID:     42261,
	},
	{
		Asset:         string(xc.SOL),
		Driver:        string(xc.DriverSolana),
		Net:           "devnet",
		URL:           "https://api.devnet.solana.com",
		ChainName:     "Solana (Devnet)",
		ExplorerURL:   "https://explorer.solana.com",
		IndexerType:   "solana",
		PollingPeriod: "2m",
		Decimals:      9,
	},
	{
		Asset:         string(xc.SUI),
		Driver:        string(xc.DriverSui),
		Net:           "devnet",
		URL:           "https://fullnode.devnet.sui.io:443",
		ChainName:     "Sui (Devnet)",
		ExplorerURL:   "https://explorer.sui.io",
		IndexerType:   "rpc",
		PollingPeriod: "2m",
		Decimals:      9,
	},
	{
		Asset:           string(xc.XPLA),
		Driver:          string(xc.DriverCosmos),
		URL:             "https://cube-rpc.xpla.dev",
		ChainName:       "XPLA (Testnet)",
		ExplorerURL:     "https://explorer.xpla.io/testnet",
		Decimals:        18,
		ChainIDStr:      "cube_47-5",
		ChainPrefix:     "xpla",
		ChainCoin:       "axpla",
		ChainCoinHDPath: 60,
	},
}
