package cmd

import (
	"fmt"
	"github.com/IBAX-io/ibax-cli/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var consoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "IBAX console, command completion",
	PreRun: loadConfigPre,
	Run: func(cmd *cobra.Command, args []string) {
		consoleStart()
	},
}

func init() {
	time.Local = time.UTC
}

var nonce int

func consoleStart() {
	if nonce == 1 {
		log.Info("Console is running")
		return
	}
	fmt.Println("\nWelcome to the IBAX console!" +
		"\nTo exit, press ctrl-d or type exit")
	line := models.NewConsole()
	if line.Nonce != nonce {
		nonce = line.Nonce
	}
	defer line.Close()
	models.NewTerminalLiner(line)
}
