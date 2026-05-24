package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for pulse.

To load completions:

Bash:
  $ source <(pulse completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ pulse completion bash > /etc/bash_completion.d/pulse
  # macOS:
  $ pulse completion bash > $(brew --prefix)/etc/bash_completion.d/pulse

Zsh:
  $ source <(pulse completion zsh)
  # To load completions for each session, execute once:
  $ pulse completion zsh > "${fpath[1]}/_pulse"

Fish:
  $ pulse completion fish | source
  # To load completions for each session, execute once:
  $ pulse completion fish > ~/.config/fish/completions/pulse.fish

PowerShell:
  PS> pulse completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
