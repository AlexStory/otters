package new

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/alexstory/otters/internal/gen"
)

var path string

func init() {
	NewCmd.Flags().StringVarP(&path, "path", "p", "", "Directory to initialize the files. Default will match the application name.")
}

var NewCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Initializes a new otters app",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := pathOrDefault(name)
		err := gen.Init(name, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func pathOrDefault(name string) string {
	if path == "" {
		return "./" + name
	}
	return path
}
