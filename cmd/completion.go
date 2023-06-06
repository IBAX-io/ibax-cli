package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `First Set application environment variables:
export PATH="$PATH:${app_path}/ibax-cli" (The current terminal is valid)

To load completions:
Bash:

$ source <(ibax-cli completion bash -d=false)

# To load completions for each session, execute once:
Linux:
  $ ibax-cli completion bash (default: /etc/bash_completion.d/ibax-cli)
MacOS:
  $ ibax-cli completion bash (default: /usr/local/etc/bash_completion.d/ibax-cli)

Zsh:

# If shell completion is not already enabled in your environment,
# you will need to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ ibax-cli completion zsh > "${fpath[1]}/_ibax-cli"

# You will need to start a new shell for this setup to take effect.

Fish:

$ ibax-cli completion fish | source

# To load completions for each session, execute once:
$ ibax-cli completion fish > ~/.config/fish/completions/ibax-cli.fish

Powershell:

PS> ibax-cli completion powershell -d=false | Out-String | Invoke-Expression

# To load completions for every new session, run:
PS> ibax-cli completion powershell  (default: ibax-cli.ps1)
# and source this file from your Powershell profile.
`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if isDefault {
				//default: linux
				path := "/etc/bash_completion.d/ibax-cli"
				if runtime.GOOS == "darwin" {
					path = "/usr/local/etc/bash_completion.d/ibax-cli"
				}
				fi, err := os.Create(path)
				if err != nil {
					log.Fatal(err)
				}
				defer fi.Close()
				cmd.Root().GenBashCompletion(fi)
			} else {
				cmd.Root().GenBashCompletion(os.Stdout)
			}
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			if isDefault {
				create, err := os.Create("ibax-cli.ps1")
				if err != nil {
					log.Fatal(err)
				}
				defer create.Close()
				cmd.Root().GenPowerShellCompletion(create)
			} else {
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		}
	},
}
var isDefault bool

func init() {
	cmdFlags := completionCmd.Flags()
	cmdFlags.BoolVarP(&isDefault, "default", "d", true, `If your completion configuration path is inconsistent with the default path, 
please disable the default configuration, 
the completion code will be output in the terminal, 
you need to write it to the specified path`)
}
