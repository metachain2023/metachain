package config

import (
	"metechain/pkg/blockchain"
	_ "metechain/pkg/crypto/sigs/secp"
	"metechain/pkg/miner"

	"github.com/spf13/viper"
)

type CfgInfo struct {
	ChainCfg    *blockchain.ChainConfig `yaml:"chainCfg"`
	SververCfg  *ServerConfig           `yaml:"sververCfg"`
	metemaskCfg *metemaskConfig         `yaml:"metemaskCfg"`
	P2PConfig   *P2PConfig              `yaml:"p2pconfig"`
	MinerConfig *miner.Config           `yaml:minerconfig`
	NetWorkType string                  `yaml:"networktype"`
}

type ServerConfig struct {
	GRpcAddress      string `yaml:"grpcaddress"`
	WebAddress       string `yaml:"webaddress"`
	InsiderpcAddress string `yaml:"insiderpcaddress"`
	GreamHost        string `yaml:"greamhost"`
	ChainServerPort  string `yaml:"chainserverport"`
}

type metemaskConfig struct {
	ChainId       string `yaml:"chainid"`
	NetworkId     string `yaml:"networkId"`
	ClinetVersion string `yaml:"clinetVersion"`
	ListenPort    string `yaml:"listenPort"`
}

type P2PConfig struct {
	AdvertiseAddr string   `yaml:"advertiseaddr"`
	Port          int      `yaml:"port"`
	JionMembers   []string `yaml:"jionmembers"`
}

// LoadConfig load configuration information
func LoadConfig() (*CfgInfo, error) {
	viper := viper.New()
	viper.SetConfigName("metechain")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg CfgInfo
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
