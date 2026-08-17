// Command shamash keeps mail local while delegating transport to cloud
// e-mail providers.
//
// It is a single binary: the submission and MX sides, the per-user index
// daemon, and the command line that operates them all live here.
//
// This is an early scaffold: the command tree exists, but none of the
// behaviour behind it is implemented.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CmdName is the name of this executable.
const CmdName = "shamash"

var rootCmd = &cobra.Command{
	Use:          CmdName,
	Short:        "Local mail storage with delegated transport",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: not implemented yet\n", CmdName)
		return err
	},
}

func main() {
	// cobra prints the error itself; SilenceUsage keeps it from dumping
	// usage on failure.
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
