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

var rootCmd = &cobra.Command{
	Use:   "traitor",
	Short: "Traitor is a privilege escalation framework for Linux",
	Long: `An extensible privilege escalation framework for Linux
                Complete documentation is available at https://github.com/liamg/traitor`,
	Run: func(cmd *cobra.Command, args []string) {

		localState := state.New()
		ctx := context.Background()
		baseLog := logger.New()
		allExploits := exploits.Get(exploits.SpeedAny)

		baseLog.Printf("Checking for opportunities...")

		for _, exploit := range allExploits {
			exploitLogger := baseLog.WithTitle(exploit.Name)
			if exploit.Vulnerability.IsVulnerable(ctx, localState, exploitLogger) {
				if disclosure, ok := exploit.Vulnerability.(exploits.Disclosure); ok {
					exploitLogger.Printf("Gathering information...")
					if err := disclosure.Disclose(ctx, localState, exploitLogger); err != nil {
						baseLog.WithTitle("error").Printf("Disclosure failed: %s", err)
					}
				}
				if sheller, ok := exploit.Vulnerability.(exploits.ShellDropper); ok {
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
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
