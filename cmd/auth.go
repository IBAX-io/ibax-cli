package cmd

import (
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
)

func loginPre(cmd *cobra.Command, args []string) {
	if models.Client != nil {
		cnf := models.Client.GetConfig()
		if cnf.Token != "" {
			models.RefreshToken()
			return
		}
		if cnf.PrivateKey == "" {
			log.Infof("Private key can't not be empty, Please set in the configuration file:%s", conf.Config.ConfigPath)
			return
		}
	}
	loadConfigPre(cmd, args)
	models.Login()
}
