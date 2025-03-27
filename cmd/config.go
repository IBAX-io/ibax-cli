package cmd

import (
	"context"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/IBAX-io/ibax-cli/models"
	"github.com/IBAX-io/ibax-cli/packages/consts"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Initial config generation",
	Run: func(cmd *cobra.Command, args []string) {
		if nonce == 1 {
			log.Info("Please exit Console")
			return
		}
		// Error omitted because we have default flag value
		configPath, _ := cmd.Flags().GetString("config")

		err := conf.FillRuntimePaths()
		if err != nil {
			log.WithError(err).Fatal("Filling config")
		}

		if configPath == "" {
			configPath = filepath.Join(conf.Config.DirPathConf.DataDir, consts.DefaultConfigFile)
		}
		err = viper.Unmarshal(&conf.Config)
		if err != nil {
			log.WithError(err).Fatal("Marshalling config to global struct variable")
		}

		err = conf.SaveConfig(configPath)
		if err != nil {
			log.WithError(err).Fatal("Saving config")
		}
		log.Infof("Config is saved to %s", configPath)
	},
}

func init() {
	cmdFlags := configCmd.Flags()
	cmdFlags.StringVar(&conf.Config.DirPathConf.DataDir, "dataDir", "", "Data directory (default cwd/data)")
	cmdFlags.StringVar(&conf.Config.DirPathConf.KeysDir, "keysDir", "", "Keys directory (default dataDir)")

	cmdFlags.StringVar(&conf.Config.Hasher, "hasher", crypto.HashAlgo_KECCAK256.String(), fmt.Sprintf("Hash Algorithm (%s | %s | %s | %s)", crypto.HashAlgo_SHA256, crypto.HashAlgo_KECCAK256, crypto.HashAlgo_SHA3_256, crypto.HashAlgo_SM3))
	cmdFlags.StringVar(&conf.Config.Cryptoer, "cryptoer", crypto.AsymAlgo_ECC_Secp256k1.String(), fmt.Sprintf("Key and Sign Algorithm (%s | %s | %s | %s)", crypto.AsymAlgo_ECC_P256, crypto.AsymAlgo_ECC_Secp256k1, crypto.AsymAlgo_ECC_P512, crypto.AsymAlgo_SM2))
	cmdFlags.Int64Var(&conf.Config.Ecosystem, "ecosystem", 1, "login ecosystem id")
	cmdFlags.StringVar(&conf.Config.RpcConnect, "connect", consts.DefaultConnect, "Send commands to node running on <connect>")
	cmdFlags.IntVar(&conf.Config.RpcPort, "port", consts.DefaultPort, "Connect to JSON-RPC on <port>")
}

// Load the configuration from file
func loadConfig(cmd *cobra.Command) {
	err := conf.LoadConfig(conf.Config.ConfigPath)
	if err != nil {
		if models.IsConsoleMode() {
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, "error", err.Error())
			cmd.SetContext(ctx)
			models.SendErrSignal(fmt.Errorf("loading config ,err: %s", err.Error()), false)
			return
		}
		log.WithError(err).Fatal("Loading config")
	}
	rpcHost := joinHost(conf.Config.RpcConnect, conf.Config.RpcPort)
	conf.UpdateSdkConfig(rpcHost)

}

func joinHost(address string, port int) (host string) {
	return fmt.Sprintf("%s:%d", address, port)
}

func loadConfigPre(cmd *cobra.Command, args []string) {
	if models.Client != nil {
		return
	}
	loadConfig(cmd)
	newClient()
}

func newClient() {
	models.Client = client.NewClient(conf.GetSdkConfig())
}
