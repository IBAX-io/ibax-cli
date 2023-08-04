package conf

import (
	"fmt"
	sdk "github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/ibax-cli/packages/consts"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// GlobalConfig is storing all startup config as global struct
type GlobalConfig struct {
	sdkConfig sdk.Config `yaml:"-"`

	PrivateKey string `json:"private_key" yaml:"private_key"` // private key. Do not use clear text. You can set the environment variable. The key controls access to your funds!
	Ecosystem  int64  `json:"ecosystem" yaml:"ecosystem"`     //Login ecosystem Id
	Cryptoer   string `json:"cryptoer" yaml:"cryptoer"`
	Hasher     string `json:"hasher" yaml:"hasher"`

	ConfigPath  string          `json:"config_path" yaml:"-"`
	RpcConnect  string          `json:"rpc_connect" yaml:"rpc_connect"`
	RpcPort     int             `json:"rpc_port" yaml:"rpc_port"`
	LinerPath   string          `json:"liner_path" yaml:"liner_path"`
	DirPathConf DirectoryConfig `json:"dir_path_conf" yaml:"dir_path_conf"`
}

type DirectoryConfig struct {
	DataDir string `json:"data_dir" yaml:"data_dir"` // application work dir (cwd by default)
	KeysDir string `json:"keys_dir" yaml:"keys_dir"` // place for private keys files: privateKey
}

var Config GlobalConfig

func LoadConfig(configPath string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\n" + `You can specify the configuration file path, use "--path" flag or make configuration, use "config" command`)
			fmt.Printf("can't not file:%s\n", configPath)
		}
		return err
	}
	configData = []byte(os.ExpandEnv(string(configData)))
	err = yaml.Unmarshal(configData, &Config)
	return err
}

// FillRuntimePaths fills paths from runtime parameters
func FillRuntimePaths() error {
	if Config.DirPathConf.DataDir == "" {
		Config.DirPathConf.DataDir = filepath.Join(consts.DefaultWorkdirName)
	}
	if Config.DirPathConf.KeysDir == "" {
		Config.DirPathConf.KeysDir = Config.DirPathConf.DataDir
	}
	if Config.LinerPath == "" {
		Config.LinerPath = Config.DirPathConf.DataDir
	}

	return nil
}

// SaveConfig save global parameters to configFile
func SaveConfig(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0775)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", dir)
		}
	}

	cf, err := os.Create(path)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Create config file failed")
		return err
	}
	defer cf.Close()

	err = yaml.NewEncoder(cf).Encode(Config)
	if err != nil {
		return err
	}

	return nil
}

func SetDefaultConfig() {
	Config.sdkConfig.EnableRpc = true
	Config.sdkConfig.JwtPrefix = "Bearer "
}

func GetSdkConfig() sdk.Config {
	return Config.sdkConfig
}

func UpdateSdkConfig(host string) {
	Config.sdkConfig.Hasher = Config.Hasher
	Config.sdkConfig.Cryptoer = Config.Cryptoer
	Config.sdkConfig.PrivateKey = Config.PrivateKey
	Config.sdkConfig.Ecosystem = Config.Ecosystem

	if host != Config.sdkConfig.ApiAddress {
		Config.sdkConfig.ApiAddress = host
	}
}
