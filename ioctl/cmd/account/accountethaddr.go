// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/ioctl/ioctlio"
	"github.com/iotexproject/iotex-core/ioctl/util"
)

// accountEthaddrCmd represents the account ethaddr command
var accountEthaddrCmd = &cobra.Command{
	Use:   "ethaddr (ALIAS|IOTEX_ADDRESS|ETH_ADDRESS)",
	Short: "Translate address between IOTX and ETH",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		err := accountEthaddr(args[0])
		return ioctlio.PrintError(err)
	},
}

type ethaddrMessage struct {
	IOAddr  string `json:"ioAddr"`
	EthAddr string `json:"ethAddr"`
}

func accountEthaddr(arg string) error {
	var ethAddress common.Address
	ioAddr, err := util.Address(arg)
	if err != nil {
		if ok := common.IsHexAddress(arg); !ok {
			return ioctlio.NewError(ioctlio.AddressError, "", err)
		}
		ethAddress = common.HexToAddress(arg)
		ioAddress, err := address.FromBytes(ethAddress.Bytes())
		if err != nil {
			return ioctlio.NewError(ioctlio.AddressError,
				fmt.Sprintf("failed to form IoTeX address from ETH address"), nil)
		}
		ioAddr = ioAddress.String()
	} else {
		ethAddress, err = util.IoAddrToEvmAddr(ioAddr)
		if err != nil {
			return ioctlio.NewError(ioctlio.AddressError, "", err)
		}
	}
	message := ethaddrMessage{IOAddr: ioAddr, EthAddr: ethAddress.String()}
	fmt.Println(message.String())
	return nil
}

func (m *ethaddrMessage) String() string {
	if ioctlio.Format == "" {
		return fmt.Sprintf("%s - %s", m.IOAddr, m.EthAddr)
	}
	return ioctlio.FormatString(ioctlio.Result, m)
}
