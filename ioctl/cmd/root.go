// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/iotexproject/iotex-core/ioctl/cmd/account"
	"github.com/iotexproject/iotex-core/ioctl/cmd/action"
	"github.com/iotexproject/iotex-core/ioctl/cmd/alias"
	"github.com/iotexproject/iotex-core/ioctl/cmd/bc"
	"github.com/iotexproject/iotex-core/ioctl/cmd/config"
	"github.com/iotexproject/iotex-core/ioctl/cmd/node"
	"github.com/iotexproject/iotex-core/ioctl/cmd/update"
	"github.com/iotexproject/iotex-core/ioctl/cmd/version"
	"github.com/iotexproject/iotex-core/ioctl/output"
)

// Multi-language support for ioctl, index 0 represents en(English).
var (
	rootCmdShorts = []string{
		"Command-line interface for IoTeX blockchain",
		"IoTeX区块链命令行工具",
	}
	rootCmdLongs = []string{
		`ioctl is a command-line interface for interacting with IoTeX blockchain.`,
		`ioctl 是用于与IoTeX区块链进行交互的命令行工具`,
	}
	flagOutputFormatUsages = []string{
		"output format",
		"指定输出格式",
	}
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ioctl",
	Short: rootCmdShorts[config.Language],
	Long:  rootCmdLongs[config.Language],
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
func init() {
	RootCmd.AddCommand(config.ConfigCmd)
	RootCmd.AddCommand(account.AccountCmd)
	RootCmd.AddCommand(alias.AliasCmd)
	RootCmd.AddCommand(action.ActionCmd)
	RootCmd.AddCommand(action.Xrc20Cmd)
	RootCmd.AddCommand(action.StakeCmd)
	RootCmd.AddCommand(bc.BCCmd)
	RootCmd.AddCommand(node.NodeCmd)
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.AddCommand(update.UpdateCmd)

	RootCmd.PersistentFlags().StringVarP(&output.Format, "output-format", "o", "",
		flagOutputFormatUsages[config.Language])
	RootCmd.HelpFunc()
}
