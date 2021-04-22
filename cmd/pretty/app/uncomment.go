package app

import (
	"github.com/27149chen/go-pretty/pkg"
	"github.com/spf13/cobra"
)

var uncommentCmd = &cobra.Command{
	Use:   "uncomment [PATH]",
	Short: "uncomment codes you do not want to expose",
	Long:  `uncomment codes you do not want to expose.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := execUncommentCmd(args[0]); err != nil {
			panic(err)
		}

		return nil
	},
}

func execUncommentCmd(root string) error {
	err := pkg.PopulateExcludes(pretty)
	if err != nil {
		return err
	}

	return pkg.Uncomment(root)
}
