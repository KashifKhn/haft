package completion

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for Haft.

Supported shells: bash, zsh, fish, powershell

To load completions:

Bash:
  # Linux:
  $ haft completion bash > /etc/bash_completion.d/haft
  
  # macOS:
  $ haft completion bash > $(brew --prefix)/etc/bash_completion.d/haft

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  # Linux:
  $ haft completion zsh > "${fpath[1]}/_haft"
  
  # macOS:
  $ haft completion zsh > $(brew --prefix)/share/zsh/site-functions/_haft

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ haft completion fish > ~/.config/fish/completions/haft.fish

PowerShell:
  PS> haft completion powershell > haft.ps1
  # and source this file from your PowerShell profile.`,
		Example: `  # Generate bash completion
  haft completion bash

  # Generate zsh completion  
  haft completion zsh

  # Generate fish completion
  haft completion fish

  # Generate powershell completion
  haft completion powershell

  # Write bash completion to file
  haft completion bash > /etc/bash_completion.d/haft`,
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(newBashCmd())
	cmd.AddCommand(newZshCmd())
	cmd.AddCommand(newFishCmd())
	cmd.AddCommand(newPowershellCmd())

	return cmd
}

func newBashCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion script",
		Long: `Generate the autocompletion script for bash.

To load completions in your current shell session:
  $ source <(haft completion bash)

To load completions for every new session:

Linux:
  $ haft completion bash > /etc/bash_completion.d/haft

macOS:
  $ haft completion bash > $(brew --prefix)/etc/bash_completion.d/haft`,
		Example: `  haft completion bash
  haft completion bash > /etc/bash_completion.d/haft`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenBashCompletion(os.Stdout)
		},
	}
}

func newZshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion script",
		Long: `Generate the autocompletion script for zsh.

If shell completion is not already enabled in your environment,
you will need to enable it by adding this to your ~/.zshrc:
  autoload -U compinit; compinit

To load completions in your current shell session:
  $ source <(haft completion zsh)

To load completions for every new session:

Linux:
  $ haft completion zsh > "${fpath[1]}/_haft"

macOS:
  $ haft completion zsh > $(brew --prefix)/share/zsh/site-functions/_haft

You will need to start a new shell for this setup to take effect.`,
		Example: `  haft completion zsh
  haft completion zsh > "${fpath[1]}/_haft"`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenZshCompletion(os.Stdout)
		},
	}
}

func newFishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fish",
		Short: "Generate fish completion script",
		Long: `Generate the autocompletion script for fish.

To load completions in your current shell session:
  $ haft completion fish | source

To load completions for every new session:
  $ haft completion fish > ~/.config/fish/completions/haft.fish

You will need to start a new shell for this setup to take effect.`,
		Example: `  haft completion fish
  haft completion fish > ~/.config/fish/completions/haft.fish`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		},
	}
}

func newPowershellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "powershell",
		Short: "Generate powershell completion script",
		Long: `Generate the autocompletion script for PowerShell.

To load completions in your current shell session:
  PS> haft completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your PowerShell profile.`,
		Example: `  haft completion powershell
  haft completion powershell > haft.ps1`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		},
	}
}
