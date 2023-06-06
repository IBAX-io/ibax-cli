package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/IBAX-io/ibax-cli/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Authorization
var (
	authStatus = &cobra.Command{
		Use:        "getAuthStatus",
		Short:      `Get Authorization Status`,
		Example:    `./ibax-cli getAuthStatus`,
		SuggestFor: []string{"getAuthStatus"},
		Args:       cobra.NoArgs,
		PreRun:     loadConfigPre,
		Run: func(cmd *cobra.Command, args []string) {
			err := cobra.NoArgs(cmd, args)
			if err != nil {
				log.Infof("no parameters required: %s", err.Error())
				return
			}
			result, err := models.Client.GetAuthStatus()
			if err != nil {
				log.Infof("Get Authorization Status Failed: %s", err.Error())
				return
			}
			if result == nil {
				log.Info("Get Authorization Status Result Empty")
				return
			}
			str, _ := json.MarshalIndent(*result, "", "    ")
			fmt.Printf("\n%+v\n", string(str))
		},
	}

	refresh = &cobra.Command{
		Use:        "refresh",
		Short:      `refresh config And re-login`,
		Example:    `./ibax-cli refresh`,
		SuggestFor: []string{"refresh"},
		Args:       cobra.NoArgs,
		Run:        refreshCmd,
	}
)

func loginPre(cmd *cobra.Command, args []string) {
	if models.Client != nil {
		cnf := models.Client.GetConfig()
		if cnf.Token != "" {
			models.RefreshToken(cmd)
			return
		}
		if cnf.PrivateKey == "" {
			err := fmt.Errorf("private key can't not be empty, Please set in the configuration file:%s", conf.Config.ConfigPath)
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, "error", err)
			cmd.SetContext(ctx)
			log.Infof(err.Error())
			return
		}
	}
	loadConfigPre(cmd, args)
	models.Login(cmd)
}

func refreshCmd(cmd *cobra.Command, args []string) {
	models.Client = nil
	path := conf.Config.ConfigPath
	rpcConnect := conf.Config.RpcConnect
	rpcPort := conf.Config.RpcPort
	dirPath := conf.Config.DirPathConf

	conf.Config = conf.GlobalConfig{
		ConfigPath:  path,
		RpcConnect:  rpcConnect,
		RpcPort:     rpcPort,
		DirPathConf: dirPath,
	}
	conf.SetDefaultConfig()
	loadConfigPre(cmd, args)
	if hasErrorContext(cmd) {
		cleanErrorContext(cmd)
		return
	}
	loginPre(cmd, args)
	if hasErrorContext(cmd) {
		cleanErrorContext(cmd)
		return
	}
	fmt.Println("\nRefresh Success!!")
}
