package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1"

func init() {
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ottr",
	Long:  "Print the version of ottr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("otter v%s\n", version)
	},
}
