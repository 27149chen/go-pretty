package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/27149chen/go-pretty/pkg"
	"github.com/27149chen/go-pretty/version"
)

var pretty string
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
		if err := run(path); err != nil {
			panic(err)
		}

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run(root string) error {
	if printVersion {
		fmt.Println(version.Version)
		return nil
	}

	err := pkg.PopulateExcludedPaths(pretty)
	if err != nil {
		return err
	}

	err = pkg.Prettify(root)
	if err != nil {
		return err
	}

	// also remove the prettyIgnore file in current directory
	dir := filepath.Dir(pretty)
	if dir != root {
		return nil
	}

	return os.Remove(pretty)
}

func init() {
	// cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&pretty, "file", "f", pkg.PrettyFile, "Name of the pretty file.")
	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version information and quit")
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
