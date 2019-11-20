// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/iotexproject/iotex-core/ioctl/output"
	"github.com/iotexproject/iotex-core/pkg/log"
)

// Directories
var (
	// ConfigDir is the directory to store config file
	ConfigDir string
	// DefaultConfigFile is the default config file name
	DefaultConfigFile string
)

// Error strings
var (
	// ErrConfigNotMatch indicates error for no config matches
	ErrConfigNotMatch = fmt.Errorf("No matching config")
	// ErrEmptyEndpoint indicates error for empty endpoint
	ErrEmptyEndpoint = fmt.Errorf("No endpoint has been set")
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get, set, or reset configuration for ioctl",
}

// Context represents the current context
type Context struct {
	AddressOrAlias string `json:"addressOrAlias" yaml:"addressOrAlias"`
}

// Config defines the config schema
type Config struct {
	Wallet         string            `json:"wallet" yaml:"wallet"`
	Endpoint       string            `json:"endpoint" yaml:"endpoint"`
	SecureConnect  bool              `json:"secureConnect" yaml:"secureConnect"`
	Aliases        map[string]string `json:"aliases" yaml:"aliases"`
	DefaultAccount Context           `json:"defaultAccount" yaml:"defaultAccount"`
	Explorer       string            `json:"explorer" yaml:"explorer"`
	Language       string            `json:"language" yaml:"language"`
}

var (
	// ReadConfig represents the current config read from local
	ReadConfig Config
	// Insecure represents the insecure connect option of grpc dial, default is false
	Insecure = false
	// Language represents the index of language of ioctl user interface, default is 0 representing English
	Language = 0
)

func init() {
	ConfigDir = os.Getenv("HOME") + "/.config/ioctl/default"
	// Create path to config directory
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		log.L().Panic(err.Error())
	}
	// Path to config file
	DefaultConfigFile = ConfigDir + "/config.default"
	// Load or reset config
	var err error
	ReadConfig, err = LoadConfig()
	if err != nil || len(ReadConfig.Wallet) == 0 {
		if err != nil && !os.IsNotExist(err) {
			log.L().Panic(err.Error()) // Config file exists but error occurs
		}
		ReadConfig.Wallet = ConfigDir
		if os.IsNotExist(err) {
			ReadConfig.SecureConnect = true
		}
		out, err := yaml.Marshal(&ReadConfig)
		if err != nil {
			log.L().Panic(err.Error())
		}
		// If default config not exist, create new default config
		if err := ioutil.WriteFile(DefaultConfigFile, out, 0600); err != nil {
			log.L().Panic(fmt.Sprintf("Failed to write to config file %s.", DefaultConfigFile))
		}
	}
	// Set language of ioctl
	Language = isSupportedLanguage(ReadConfig.Language)
	if Language == -1 {
		Language = 0
		message := output.StringMessage(fmt.Sprintf("Language %s is not supported, English instead."))
		fmt.Println(message.Warn())
	}
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configResetCmd)
}

// LoadConfig loads config file in yaml format
func LoadConfig() (Config, error) {
	ReadConfig := Config{
		Aliases: make(map[string]string),
	}
	in, err := ioutil.ReadFile(DefaultConfigFile)
	if err == nil {
		if err := yaml.Unmarshal(in, &ReadConfig); err != nil {
			return ReadConfig, err
		}
	}
	return ReadConfig, err
}
