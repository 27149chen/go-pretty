package app

import (
	"github.com/27149chen/go-pretty/pkg/pretty"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment [PATH]",
	Short: "comment codes you do not want to expose",
	Long:  `comment codes you do not want to expose instead of removing them`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := execCommentCmd(args[0]); err != nil {
			panic(err)
		}

		return nil
	},
}

func execCommentCmd(root string) error {
	err := pretty.PopulateExcludes(prettyFile)
	if err != nil {
		return err
	}

	return pretty.Comment(root)
}
