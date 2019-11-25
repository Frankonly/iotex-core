// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package alias

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/iotexproject/iotex-core/ioctl/cmd/config"
	"github.com/iotexproject/iotex-core/ioctl/ioctlio"
)

// aliasImportCmd represents the alias import command
var aliasImportCmd = &cobra.Command{
	Use:   "import 'DATA'",
	Short: "Import aliases",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		err := aliasImport(cmd, args)
		return ioctlio.PrintError(err)
	},
}

type importMessage struct {
	ImportedNumber int     `json:"importedNumber"`
	TotalNumber    int     `json:"totalNumber"`
	Imported       []alias `json:"imported"`
	Unimported     []alias `json:"unimported"`
}

func init() {
	aliasImportCmd.Flags().StringVarP(&format,
		"format=", "f", "json", "set format: json/yaml")
	aliasImportCmd.Flags().BoolVarP(&forceImport,
		"force-import", "F", false, "override existing aliases")
}

func aliasImport(cmd *cobra.Command, args []string) error {
	var importedAliases aliases
	switch format {
	default:
		return ioctlio.NewError(ioctlio.FlagError, fmt.Sprintf("invalid format flag %s", format), nil)
	case "json":
		if err := json.Unmarshal([]byte(args[0]), &importedAliases); err != nil {
			return ioctlio.NewError(ioctlio.SerializationError, "failed to unmarshal imported aliases", err)
		}
	case "yaml":
		if err := yaml.Unmarshal([]byte(args[0]), &importedAliases); err != nil {
			return ioctlio.NewError(ioctlio.SerializationError, "failed to unmarshal imported aliases", err)
		}
	}
	aliases := GetAliasMap()
	message := importMessage{TotalNumber: len(importedAliases.Aliases), ImportedNumber: 0}
	for _, importedAlias := range importedAliases.Aliases {
		if !forceImport && config.ReadConfig.Aliases[importedAlias.Name] != "" {
			message.Unimported = append(message.Unimported, importedAlias)
			continue
		}
		for aliases[importedAlias.Address] != "" {
			delete(config.ReadConfig.Aliases, aliases[importedAlias.Address])
			aliases = GetAliasMap()
		}
		config.ReadConfig.Aliases[importedAlias.Name] = importedAlias.Address
		message.Imported = append(message.Imported, importedAlias)
		message.ImportedNumber++
	}
	out, err := yaml.Marshal(&config.ReadConfig)
	if err != nil {
		return ioctlio.NewError(ioctlio.SerializationError, "failed to marshal config", err)
	}
	if err := ioutil.WriteFile(config.DefaultConfigFile, out, 0600); err != nil {
		return ioctlio.NewError(ioctlio.WriteFileError,
			fmt.Sprintf("failed to write to config file %s", config.DefaultConfigFile), err)
	}
	fmt.Println(message.String())
	return nil
}

func (m *importMessage) String() string {
	if ioctlio.Format == "" {
		line := fmt.Sprintf("%d/%d aliases imported\nExisted aliases:", m.ImportedNumber, m.TotalNumber)
		for _, alias := range m.Unimported {
			line += fmt.Sprint(" " + alias.Name)
		}
		return line
	}
	return ioctlio.FormatString(ioctlio.Result, m)
}
