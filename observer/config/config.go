package config

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ClientConfiguration
type Configuration struct {
	ChainID      string `json:"chain_id" mapstructure:"chain_id" `
	ChainHost    string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC     string `json:"chain_rpc" mapstructure:"chain_rpc"`
	SignerName   string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd string `json:"signer_passwd" mapstructure:"signer_passwd"`
	NodeRPC      string `json:"node_rpc" mapstructure:"node_rpc"`
}

// LoadConfig read the tss-proc configuration from the given file
func LoadConfig(homePath string, file string) (*Configuration, error) {
	var cfg Configuration
	viper.AddConfigPath(homePath)
	viper.AddConfigPath(filepath.Dir(file))
	viper.SetConfigName(strings.TrimRight(path.Base(file), ".json"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fail to read from config file: %w", err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("fail to unmarshal: %w", err)
	}

	return &cfg, nil
}
