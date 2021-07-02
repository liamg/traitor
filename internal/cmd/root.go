package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/liamg/traitor/internal/version"

	"github.com/liamg/traitor/pkg/logger"
	"github.com/liamg/traitor/pkg/state"

	"github.com/liamg/traitor/pkg/exploits"

	"github.com/spf13/cobra"
)

var runAnyExploit bool
var exploitName string
var promptForPassword bool
var skipExploits []string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&runAnyExploit, "any", "a", runAnyExploit, "Attempt to exploit a vulnerability as soon as it is detected. Provides a shell where possible.")
	rootCmd.PersistentFlags().BoolVarP(&promptForPassword, "with-password", "p", promptForPassword, "Prompt for the user password, if you know it. Can provide more GTFOBins possibilities via sudo.")
	rootCmd.PersistentFlags().StringVarP(&exploitName, "exploit", "e", exploitName, "Run the specified exploit, if the system is found to be vulnerable. Provides a shell where possible.")
	rootCmd.PersistentFlags().StringSliceVarP(&skipExploits, "skip", "k", skipExploits, "Exploit(s) to skip - specify multiple times to skip multiple exploits.")
}

var rootCmd = &cobra.Command{
	Use:   "traitor",
	Short: "Traitor is a privilege escalation framework for Linux",
	Long: `An extensible privilege escalation framework for Linux
                Complete documentation is available at https://github.com/liamg/traitor`,
	Args: cobra.ExactArgs(0),
	PreRun: func(_ *cobra.Command, args []string) {
		fmt.Printf("\x1b[34m"+`

▀█▀ █▀█ ▄▀█ █ ▀█▀ █▀█ █▀█
░█░ █▀▄ █▀█ █ ░█░ █▄█ █▀▄ %s
`+"\x1b[31mhttps://github.com/liamg/traitor\n\n", version.Version)
	},
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		baseLog := logger.New()

		if user, err := user.Current(); err == nil && user.Uid == "0" {
			baseLog.Printf("You are already root.")
			return
		}

		baseLog.Printf("Assessing machine state...")
		localState := state.New()
		localState.HasPassword = promptForPassword
		localState.Assess()

		baseLog.Printf("Checking for opportunities...")
		allExploits := exploits.Get(exploits.SpeedAny)
		var found bool
		var vulnFound bool
		for _, exploit := range allExploits {
			if shouldSkip(exploit.Name) {
				continue
			}
			if exploitName == "" || exploitName == exploit.Name {
				found = true
				exploitLogger := baseLog.WithTitle(exploit.Name)
				if exploit.Vulnerability.IsVulnerable(ctx, localState, exploitLogger) {
					vulnFound = true
					if disclosure, ok := exploit.Vulnerability.(exploits.Disclosure); ok {
						exploitLogger.Printf("Gathering information...")
						if err := disclosure.Disclose(ctx, localState, exploitLogger); err != nil {
							baseLog.WithTitle("error").Printf("Disclosure failed: %s", err)
						}
					}
					if sheller, ok := exploit.Vulnerability.(exploits.ShellDropper); ok {
						if runAnyExploit || exploitName == exploit.Name {
							exploitLogger.Printf("Opportunity found, trying to exploit it...")
							if err := sheller.Shell(ctx, localState, exploitLogger); err != nil {
								baseLog.WithTitle("error").Printf("Exploit failed: %s", err)
								baseLog.Printf("Continuing to look for opportunities")
								vulnFound = false
								continue
							}
							exploitLogger.Printf("Session complete.")
							baseLog.Printf("Done.")
							return
						}
						exploitLogger.Printf("System is vulnerable! Run again with '--exploit %s' to exploit it.", exploit.Name)
					} else if exploitName != "" {
						exploitLogger.Printf("No local exploit available for '%s'", exploit.Name)
					}
				} else if exploitName != "" {
					exploitLogger.Printf("System is not vulnerable to '%s' - cannot exploit.", exploit.Name)
				}
			}
		}
		if exploitName != "" && !found {
			baseLog.Printf("No exploit found for '%s'", exploitName)
		} else if !vulnFound {
			baseLog.Printf("Nothing found to exploit.")
		}
	},
}

func shouldSkip(id string) bool {
	for _, skip := range skipExploits {
		if skip == id {
			return true
		}
	}
	return false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
