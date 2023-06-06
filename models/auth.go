package models

import (
	"context"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/ibax-cli/conf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

func Login(cmd *cobra.Command) {
	err := Client.AutoLogin()
	if err != nil {
		if IsConsoleMode() {
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, "error", err)
			cmd.SetContext(ctx)
			SendErrSignal(fmt.Errorf("[login] Authorization failed:%s", err.Error()), false)
			return
		}
		log.Fatalf("[login] Authorization failed: %s", err.Error())
		return
	}
}

func RefreshToken(cmd *cobra.Command) {
	cfg := Client.GetConfig()
	if time.Unix(cfg.TokenExpireTime, 0).Sub(time.Now()) < time.Minute*10 {
		cfg.Token = ""
		Client.SetConfig(cfg)
		err := Client.AutoLogin()
		if err != nil {
			if IsConsoleMode() {
				ctx := cmd.Context()
				ctx = context.WithValue(ctx, "error", err)
				cmd.SetContext(ctx)
				SendErrSignal(fmt.Errorf("[refresh token] Authorization failed:%s", err.Error()), false)
				return
			}
			log.Fatalf("[refresh token] Authorization failed: %s", err.Error())
			return
		}
	}
}

func NewClient() {
	Client = client.NewClient(conf.Config.SdkConfig)
}
