package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/27149chen/go-pretty/pkg/config"
	"github.com/27149chen/go-pretty/pkg/pretty"
	"github.com/27149chen/go-pretty/version"
)

var prettyFile string
var printVersion bool

var rootCmd = &cobra.Command{
	Use:   "pretty [PATH]",
	Short: "Prettify your project by removing things you do not want to expose",
	Long: `Prettify your project by removing things you do not want to expose.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		if err := execRootCmd(path); err != nil {
			panic(err)
		}

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func execRootCmd(root string) error {
	if printVersion {
		fmt.Println(version.Version)
		return nil
	}

	err := pretty.PopulateExcludes(prettyFile)
	if err != nil {
		return err
	}

	err = pretty.Prettify(root)
	if err != nil {
		return err
	}

	// also remove the prettyIgnore file in current directory
	dir := filepath.Dir(prettyFile)
	if dir == root {
		if err := os.Remove(prettyFile); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
	}
	
	return pretty.RemoveEmptyDir(root)
}

func init() {
	// cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&prettyFile, "file", "f", config.DefaultPrettyFile, "Name of the pretty file.")
	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version information and quit")

	rootCmd.AddCommand(commentCmd)
	rootCmd.AddCommand(uncommentCmd)
}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		cobra.CheckErr(err)
//
//		// Search config in home directory with name ".go-pretty" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".go-pretty")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
//	}
//}
