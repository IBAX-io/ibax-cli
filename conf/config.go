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
	SdkConfig   sdk.Config      `json:"sdk_config" yaml:"sdk_config"`
	ConfigPath  string          `json:"config_path" yaml:"-"`
	RpcConnect  string          `json:"rpc_connect" yaml:"-"`
	RpcPort     int             `json:"rpc_port" yaml:"-"`
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
	Config.SdkConfig.EnableRpc = true
	Config.SdkConfig.JwtPrefix = "Bearer "
}
