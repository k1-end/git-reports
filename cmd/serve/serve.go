package serve

import (
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the report",
	Long:  "Serve the report on your local network and view it in your browser",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
