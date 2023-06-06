package models

import (
	"fmt"
	"github.com/spf13/cobra"
)

var globalCmd cobra.Command

func InitGlobalCmd(cmd *cobra.Command) {
	if cmd != nil {
		globalCmd = *cmd
	} else {
		panic(fmt.Errorf("init global cmd failed"))
	}
}
