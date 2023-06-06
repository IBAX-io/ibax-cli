package cmd

import (
	"fmt"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/IBAX-io/ibax-cli/models"
	"github.com/IBAX-io/ibax-cli/packages/consts"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ibax-cli",
	Short: fmt.Sprintf("IBAX Core RPC Client Version: %s", consts.Version()),
}

var (
	buildBranch = ""
	buildDate   = ""
	commitHash  = ""
)

func init() {
	accountCmd.AddCommand(
		accountNewCmd,
		accountListCmd,
		keyIdToAddressCmd,
		addressToKeyIdCmd,
		publicKeyToAddressCmd,
		infoCmd,
	)
	for _, subCommand := range accountCmd.Commands() {
		for k, v := range subCommand.SuggestFor {
			subCommand.SuggestFor[k] = accountCmd.Use + " " + v
			models.AddWordsCompletions(subCommand.SuggestFor)
		}
	}

	rootCmd.AddCommand(
		configCmd,
		versionCmd,
		completionCmd,
		consoleCmd,
		accountCmd,
	)

	initCmdList()
	for _, c := range cmdList {
		models.AddWordsCompletions(c.SuggestFor)
		rootCmd.AddCommand(c)
	}

	consts.BuildInfo = func() string {
		if buildBranch == "" {
			return fmt.Sprintf("branch.%s commit.%s time.%s", "unknown", "unknown", "unknown")
		}
		return fmt.Sprintf("branch.%s commit.%s time.%s", buildBranch, commitHash, buildDate)
	}()

	cmdFlags := rootCmd.PersistentFlags()
	// This flags are visible for all child commands
	cmdFlags.StringVar(&conf.Config.ConfigPath, "path", defaultConfigPath(), "filepath to config.yml")
	cmdFlags.StringVar(&conf.Config.RpcConnect, "rcpConnect", consts.DefaultConnect, "Send commands to node running on <ip>")
	cmdFlags.IntVar(&conf.Config.RpcPort, "rpcPort", consts.DefaultPort, "Connect to JSON-RPC on <port>")

	conf.SetDefaultConfig()
	viper.BindPFlags(cmdFlags)
	models.InitGlobalCmd(rootCmd)
}

func defaultConfigPath() string {
	return filepath.Join("data", "config.yml")
}

// Execute executes rootCmd command.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Executing root command")
	}
}
