package cmd

import (
	"fmt"
	"github.com/IBAX-io/ibax-cli/packages/consts"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(consts.Version())
	},
}
