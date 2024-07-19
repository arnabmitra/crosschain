package main

import (
	xc "github.com/cordialsys/crosschain"
	"github.com/cordialsys/crosschain/cmd/xc/setup"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CmdXc() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "xc",
		Short:        "Manually interact with blockchains",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			args, err := setup.RpcArgsFromCmd(cmd)
			if err != nil {
				return err
			}
			setup.ConfigureLogger(args)

			xcFactory, err := setup.LoadFactory(args)
			if err != nil {
				return err
			}
			chainConfig, err := setup.LoadChain(xcFactory, args.Chain)
			if err != nil {
				return err
			}
			setup.OverrideChainSettings(chainConfig, args)

			ctx := setup.CreateContext(xcFactory, chainConfig)

			logrus.WithFields(logrus.Fields{
				"rpc":     chainConfig.GetAllClients()[0].URL,
				"network": chainConfig.GetAllClients()[0].Network,
				"chain":   chainConfig.Chain,
			}).Info("chain")
			cmd.SetContext(ctx)
			return nil
		},
	}
	setup.AddRpcArgs(cmd)

	cmd.AddCommand(CmdRpcBalance())
	cmd.AddCommand(CmdTxInput())
	cmd.AddCommand(CmdTxInfo())
	cmd.AddCommand(CmdTxTransfer())
	cmd.AddCommand(CmdAddress())
	cmd.AddCommand(CmdChains())

	return cmd
}

func assetConfig(chain *xc.ChainConfig, contractMaybe string, decimals int32) xc.ITask {
	if contractMaybe != "" {
		token := xc.TokenAssetConfig{
			Contract:    contractMaybe,
			Chain:       chain.Chain,
			ChainConfig: chain,
			Decimals:    decimals,
		}
		return &token
	} else {
		return chain
	}
}

func main() {
	rootCmd := CmdXc()
	_ = rootCmd.Execute()
}
