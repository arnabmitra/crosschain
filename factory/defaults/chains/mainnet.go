package chains

import (
	xc "github.com/jumpcrypto/crosschain"
)

func init() {
	for _, chain := range Mainnet {
		if chain.Net == "" {
			chain.Net = "mainnet"
		}

		// default to using xc client
		if chain.URL == "" && len(chain.Clients) == 0 {
			chain.Clients = append(chain.Clients, &xc.ClientConfig{
				Driver: string(xc.DriverCrosschain),
				URL:    "https://crosschain.cordialapis.com",
			})
		}
	}
}

var Mainnet = []*xc.NativeAssetConfig{
	{
		Asset:         string(xc.ACA),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://eth-rpc-acala.aca-api.network",
		ChainName:     "Acala",
		ExplorerURL:   "https://acala.subscan.io",
		IndexerUrl:    "",
		IndexerType:   "rpc",
		PollingPeriod: "15m",
		Decimals:      12,
		ChainID:       787,
	},
	{
		Asset:         string(xc.APTOS),
		Driver:        string(xc.DriverAptos),
		URL:           "https://aptos-mainnet-rpc.allthatnode.com",
		ChainName:     "Aptos",
		ExplorerURL:   "https://explorer.aptoslabs.com/",
		IndexerUrl:    "https://indexer.mainnet.aptoslabs.com",
		IndexerType:   "aptos",
		PollingPeriod: "8m",
		Decimals:      8,
		ChainID:       1,
	},
	{
		Asset:                string(xc.ATOM),
		Driver:               string(xc.DriverCosmos),
		URL:                  "https://cosmos-rpc.publicnode.com",
		ChainName:            "Cosmos",
		ExplorerURL:          "https://atomscan.com",
		IndexerType:          "cosmos",
		PollingPeriod:        "5m",
		Decimals:             6,
		ChainIDStr:           "cosmoshub-4",
		ChainPrefix:          "cosmos",
		ChainCoin:            "uatom",
		ChainCoinHDPath:      118,
		ChainGasPriceDefault: 0.100000,
	},
	{
		Asset:         string(xc.AVAX),
		Driver:        string(xc.DriverEVM),
		URL:           "https://api.avax.network/ext/bc/C/rpc",
		ChainName:     "Avalanche C-Chain",
		ExplorerURL:   "https://snowtrace.io",
		IndexerType:   "covalent",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       43114,
		Clients: []*xc.ClientConfig{
			{
				Driver: "crosschain",
				URL:    "https://crosschain.cordialapis.com",
			},
		},
	},
	{
		Asset:              string(xc.ArbETH),
		Driver:             string(xc.DriverEVM),
		URL:                "https://arb1.arbitrum.io/rpc",
		ChainName:          "Arbitrum",
		ExplorerURL:        "https://arbiscan.io",
		IndexerType:        "rpc",
		PollingPeriod:      "10m",
		Decimals:           18,
		ChainID:            42161,
		ChainGasMultiplier: 0.05,
	},
	{
		Asset:         string(xc.AurETH),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://mainnet.aurora.dev",
		ChainName:     "Aurora",
		ExplorerURL:   "https://aurorascan.dev",
		IndexerType:   "rpc",
		PollingPeriod: "15m",
		Decimals:      18,
		ChainID:       1313161554,
	},
	{
		Asset:         string(xc.BCH),
		Driver:        "bitcoin",
		URL:           "https://api.blockchair.com/bitcoin-cash",
		ChainName:     "Bitcoin Cash",
		ExplorerURL:   "https://blockchair.com/bitcoin-cash",
		IndexerUrl:    "https://api.blockchair.com/bitcoin-cash",
		IndexerType:   "blockchair",
		PollingPeriod: "10m",
		Decimals:      8,
		Provider:      "blockchair",
		Auth:          "ENV:BLOCKCHAIN_API_KEY",
	},
	{
		Asset:         string(xc.BNB),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://bsc-dataseed.binance.org",
		ChainName:     "Binance Smart Chain",
		ExplorerURL:   "https://bscscan.com",
		IndexerType:   "covalent",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       56,
	},
	{
		Asset:         string(xc.BTC),
		Driver:        string(xc.DriverBitcoin),
		URL:           "https://api.blockchair.com/bitcoin",
		ChainName:     "Bitcoin",
		ExplorerURL:   "https://blockchair.com/bitcoin",
		IndexerUrl:    "https://api.blockchair.com/bitcoin",
		IndexerType:   "blockchair",
		PollingPeriod: "10m",
		Decimals:      8,
		Provider:      "blockchair",
		Auth:          "ENV:BLOCKCHAIN_API_KEY",
	},
	{
		Asset:         string(xc.CELO),
		Driver:        string(xc.DriverEVM),
		URL:           "https://forno.celo.org",
		ChainName:     "Celo",
		ExplorerURL:   "https://explorer.celo.org",
		IndexerUrl:    "https://explorer.celo.org/mainnet",
		IndexerType:   "blockscout",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       42220,
	},
	{
		Asset:         string(xc.CHZ2),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://rpc.chiliz.com",
		ChainName:     "Chiliz 2.0",
		ExplorerURL:   "https://scan.chiliz.com",
		IndexerUrl:    "https://scan.chiliz.com",
		IndexerType:   "blockscout",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       88888,
	},
	{
		Asset:         string(xc.CHZ),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "",
		ChainName:     "Chiliz",
		ExplorerURL:   "https://explorer.chiliz.com",
		IndexerUrl:    "https://explorer.chiliz.com",
		IndexerType:   "blockscout",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       99999,
		// CHZ legacy chain has private RPC which you must provide an override for
		Disabled: true,
	},
	{
		Asset:         string(xc.DOGE),
		Driver:        string(xc.DriverBitcoin),
		URL:           "https://api.blockchair.com/dogecoin",
		ChainName:     "Dogecoin",
		ExplorerURL:   "https://blockchair.com/dogecoin",
		IndexerUrl:    "https://api.blockchair.com/dogecoin",
		IndexerType:   "blockchair",
		PollingPeriod: "10m",
		Decimals:      8,
		Provider:      "blockchair",
		Auth:          "ENV:BLOCKCHAIN_API_KEY",
	},
	{
		Asset:         string(xc.ETC),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://www.ethercluster.com/etc",
		ChainName:     "Ethereum Classic",
		ExplorerURL:   "https://blockscout.com/etc/mainnet",
		IndexerUrl:    "https://blockscout.com/etc/mainnet",
		IndexerType:   "blockscout",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       61,
	},
	{
		Asset:         string(xc.ETH),
		Driver:        string(xc.DriverEVM),
		URL:           "https://ethereum.publicnode.com",
		ChainName:     "Ethereum",
		ExplorerURL:   "https://etherscan.io",
		IndexerType:   "covalent",
		PollingPeriod: "3m",
		Decimals:      18,
		ChainID:       1,
	},
	{
		Asset:         string(xc.ETHW),
		Driver:        string(xc.DriverEVM),
		URL:           "http://mainnet.ethereumpow.org:80",
		ChainName:     "EthereumPOW",
		ExplorerURL:   "https://etherscan.io",
		IndexerType:   "rpc",
		PollingPeriod: "5m",
		Decimals:      18,
		ChainID:       10001,
	},
	{
		Asset:         string(xc.FTM),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://rpc.ftm.tools",
		ChainName:     "Fantom",
		ExplorerURL:   "https://ftmscan.com",
		IndexerType:   "covalent",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       250,
	},
	{
		Asset:                string(xc.INJ),
		Driver:               string(xc.DriverCosmos),
		URL:                  "https://injective-rpc.polkachu.com",
		ChainName:            "Injective",
		ExplorerURL:          "https://explorer.injective.network",
		IndexerType:          "cosmos",
		PollingPeriod:        "5m",
		Decimals:             18,
		ChainIDStr:           "injective-1",
		ChainPrefix:          "inj",
		ChainCoin:            "inj",
		ChainCoinHDPath:      60,
		ChainGasPriceDefault: 500000000.000000,
	},
	{
		Asset:         string(xc.KAR),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://eth-rpc-karura.aca-api.network",
		ChainName:     "Karura",
		ExplorerURL:   "https://karura.subscan.io",
		IndexerType:   "rpc",
		PollingPeriod: "15m",
		Decimals:      12,
		ChainID:       686,
	},
	{
		Asset:         string(xc.KLAY),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://public-node-api.klaytnapi.com/v1/cypress",
		ChainName:     "Klaytn",
		ExplorerURL:   "https://scope.klaytn.com",
		IndexerType:   "rpc",
		PollingPeriod: "15m",
		Decimals:      18,
		ChainID:       8217,
	},
	{
		Asset:         string(xc.LTC),
		Driver:        string(xc.DriverBitcoin),
		URL:           "https://api.blockchair.com/litecoin",
		ChainName:     "Litecoin",
		ExplorerURL:   "https://blockchair.com/litecoin",
		IndexerUrl:    "https://api.blockchair.com/litecoin",
		IndexerType:   "blockchair",
		PollingPeriod: "10m",
		Decimals:      8,
		Provider:      "blockchair",
		Auth:          "ENV:BLOCKCHAIN_API_KEY",
	},
	{
		Asset:              string(xc.LUNA),
		Driver:             string(xc.DriverCosmos),
		URL:                "https://terra-rpc.publicnode.com",
		ChainName:          "Terra",
		ExplorerURL:        "https://finder.terra.money",
		IndexerType:        "cosmos",
		PollingPeriod:      "5m",
		Decimals:           6,
		ChainGasMultiplier: 12.00,
		ChainIDStr:         "phoenix-1",
		ChainPrefix:        "terra",
		ChainCoin:          "uluna",
		ChainCoinHDPath:    330,
	},
	{
		Asset:                string(xc.LUNC),
		Driver:               string(xc.DriverCosmos),
		URL:                  "https://terra-classic-rpc.publicnode.com",
		ChainName:            "Terra Classic",
		ExplorerURL:          "https://finder.terra.money/classic",
		IndexerType:          "cosmos",
		PollingPeriod:        "5m",
		Decimals:             6,
		ChainGasMultiplier:   2.00,
		ChainIDStr:           "columbus-5",
		ChainPrefix:          "terra",
		ChainCoin:            "uluna",
		GasCoin:              "uusd",
		ChainCoinHDPath:      330,
		ChainGasPriceDefault: 5.000000,
		// Terra classic has a 0.5% tax on all bank transfers
		ChainTransferTax: 0.005,
	},
	{
		Asset:              string(xc.MATIC),
		Driver:             string(xc.DriverEVM),
		URL:                "https://polygon-rpc.com",
		ChainName:          "Polygon",
		ExplorerURL:        "https://polygonscan.com",
		IndexerType:        "covalent",
		PollingPeriod:      "6m",
		Decimals:           18,
		ChainID:            137,
		ChainGasMultiplier: 4.00,
		ChainGasTip:        20,
	},
	{
		Asset:         string(xc.OAS),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://rpc.mainnet.oasys.games",
		ChainName:     "Oasys",
		ExplorerURL:   "https://explorer.oasys.games",
		IndexerUrl:    "https://scan.oasys.games",
		IndexerType:   "blockscout",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       248,
	},
	{
		Asset:         string(xc.OptETH),
		Driver:        string(xc.DriverEVM),
		URL:           "https://mainnet.optimism.io",
		ChainName:     "Optimism",
		ExplorerURL:   "https://optimistic.etherscan.io",
		IndexerType:   "rpc",
		PollingPeriod: "15m",
		Decimals:      18,
		ChainID:       10,
	},
	{
		Asset:         string(xc.ROSE),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://emerald.oasis.dev",
		ChainName:     "Oasis",
		ExplorerURL:   "https://explorer.emerald.oasis.dev",
		IndexerType:   "covalent",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       42262,
	},
	{
		Asset:         string(xc.SOL),
		Driver:        string(xc.DriverSolana),
		URL:           "https://api.mainnet-beta.solana.com",
		ChainName:     "Solana",
		ExplorerURL:   "https://explorer.solana.com",
		IndexerType:   "solana",
		PollingPeriod: "3m",
		Decimals:      9,
		Clients: []*xc.ClientConfig{
			{
				Driver: "crosschain",
				URL:    "https://crosschain.cordialapis.com",
			},
		},
	},
	{
		Asset:         string(xc.SUI),
		Driver:        string(xc.DriverSui),
		URL:           "https://fullnode.mainnet.sui.io:443",
		ChainName:     "Sui",
		ExplorerURL:   "https://explorer.sui.io",
		IndexerType:   "rpc",
		PollingPeriod: "3m",
		Decimals:      9,
	},
	{
		Asset:         string(xc.XDC),
		Driver:        string(xc.DriverEVMLegacy),
		URL:           "https://rpc.xdcrpc.com",
		ChainName:     "XinFin",
		ExplorerURL:   "https://explorer.xinfin.network/",
		IndexerUrl:    "https://xdc.blocksscan.io",
		IndexerType:   "blocksscan",
		PollingPeriod: "6m",
		Decimals:      18,
		ChainID:       50,
	},
	{
		Asset:                string(xc.XPLA),
		Driver:               string(xc.DriverCosmos),
		URL:                  "https://dimension-rpc.xpla.dev",
		ChainName:            "XPLA Chain",
		ExplorerURL:          "https://explorer.xpla.io/mainnet",
		IndexerType:          "cosmos",
		PollingPeriod:        "5m",
		Decimals:             18,
		ChainIDStr:           "dimension_37-1",
		ChainPrefix:          "xpla",
		ChainCoin:            "axpla",
		ChainCoinHDPath:      60,
		ChainGasPriceDefault: 0.1,
	},

	{
		Asset:                string(xc.TIA),
		Driver:               string(xc.DriverCosmos),
		URL:                  "",   // TODO update when celestia mainnet is ready
		ChainIDStr:           "",   // UPDATE
		Disabled:             true, // UPDATE
		ChainName:            "Celestia",
		ExplorerURL:          "",
		IndexerType:          "cosmos",
		PollingPeriod:        "15m",
		Decimals:             6,
		ChainPrefix:          "celestia",
		ChainCoin:            "utia",
		ChainCoinHDPath:      118,
		ChainGasMultiplier:   12.0,
		ChainGasPriceDefault: 0.1,
	},
	{
		Asset:                string(xc.SEI),
		Driver:               string(xc.DriverCosmos),
		URL:                  "",   // TODO update for mainnet when ready
		ChainIDStr:           "",   // UPDATE
		Disabled:             true, // UPDATE
		ChainPrefix:          "sei",
		ChainCoin:            "usei",
		ChainCoinHDPath:      118,
		ChainName:            "Sei",
		ChainGasMultiplier:   12.0,
		ChainGasPriceDefault: 0.1,
		ExplorerURL:          "https://sei.explorers.guru/",
		IndexerType:          "cosmos",
		PollingPeriod:        "15m",
		Decimals:             6,
	},
}
