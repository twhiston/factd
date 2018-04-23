package cmd

import (
	"github.com/spf13/cobra"
	"github.com/twhiston/factd/lib/common/logging"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs factd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		factd := setupFactD()
		factd.Collect()
		logging.HandleError(factd.Print())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
