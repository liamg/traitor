package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/liamg/traitor/internal/logger"
	"github.com/liamg/traitor/pkg/state"

	"github.com/liamg/traitor/pkg/exploits"

	"github.com/spf13/cobra"
)

var runAnyExploit bool
var exploitName string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&runAnyExploit, "exploit-any", "a", runAnyExploit, "Attempt to exploit a vulnerability as soon as it is detected. Provides a shell where possible.")
	rootCmd.PersistentFlags().StringVarP(&exploitName, "exploit", "e", exploitName, "Run the specified exploit, if the system is found to be vulnerable. Provides a shell where possible.")
}

var rootCmd = &cobra.Command{
	Use:   "traitor",
	Short: "Traitor is a privilege escalation framework for Linux",
	Long: `An extensible privilege escalation framework for Linux
                Complete documentation is available at https://github.com/liamg/traitor`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		baseLog := logger.New()

		baseLog.Printf("Assessing machine state...")
		localState := state.New()
		localState.Assess()

		baseLog.Printf("Checking for opportunities...")
		allExploits := exploits.Get(exploits.SpeedAny)
		var found bool
		for _, exploit := range allExploits {
			if exploitName == "" || exploitName == exploit.Name {
				found = true
				exploitLogger := baseLog.WithTitle(exploit.Name)
				if exploit.Vulnerability.IsVulnerable(ctx, localState, exploitLogger) {
					if disclosure, ok := exploit.Vulnerability.(exploits.Disclosure); ok {
						exploitLogger.Printf("Gathering information...")
						if err := disclosure.Disclose(ctx, localState, exploitLogger); err != nil {
							baseLog.WithTitle("error").Printf("Disclosure failed: %s", err)
						}
					}
					if sheller, ok := exploit.Vulnerability.(exploits.ShellDropper); ok {
						if runAnyExploit {
							exploitLogger.Printf("System is vulnerable, starting exploit...")
							if err := sheller.Shell(ctx, localState, exploitLogger); err != nil {
								baseLog.WithTitle("error").Printf("Exploit failed: %s", err)
								baseLog.Printf("Continuing to look for opportunities")
								continue
							}
							exploitLogger.Printf("Exploit successful.")
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
		if !found {
			baseLog.Printf("No exploit found for '%s'", exploitName)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
