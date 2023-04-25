package models

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/ibax-cli/conf"
	log "github.com/sirupsen/logrus"
	"time"
)

func Login() {
	err := Client.AutoLogin()
	if err != nil {
		log.Fatalf("[login] Authorization failed: %s", err.Error())
		return
	}
}

func RefreshToken() {
	cfg := Client.GetConfig()
	if time.Unix(cfg.TokenExpireTime, 0).Sub(time.Now()) < time.Minute*10 {
		cfg.Token = ""
		Client.SetConfig(cfg)
		err := Client.AutoLogin()
		if err != nil {
			log.Fatalf("[refresh token] Authorization failed: %s", err.Error())
			return
		}
	}
}

func NewClient() {
	Client = client.NewClient(conf.Config.SdkConfig)
}
