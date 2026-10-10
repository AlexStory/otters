package gen

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/alexstory/otters/internal/gen"
)

var (
	genPath   string
	genPrefix string
)

var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Generate a sub-app and mount it into main.go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		projectDir := genPath
		if projectDir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			projectDir = cwd
		}

		if err := gen.GenApp(projectDir, name, genPrefix); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	AppCmd.Flags().StringVarP(&genPath, "path", "p", "", "path to the project directory")
	AppCmd.Flags().StringVar(&genPrefix, "prefix", "", "Mount prefix (defaults to /<name>")
	GenCmd.AddCommand(AppCmd)
}
