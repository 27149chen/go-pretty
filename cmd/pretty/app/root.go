package app

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var pretty string

var rootCmd = &cobra.Command{
	Use:   "pretty PATH",
	Short: "Prettify your project by removing things you do not want to expose",
	Long: `Prettify your project by removing things you do not want to expose.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run(args[0]); err != nil {
			panic(err)
		}

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run(root string) error {
	err := populateExcludedPaths(pretty)
	if err != nil {
		return err
	}

	err = prettify(root)
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

	rootCmd.PersistentFlags().StringVarP(&pretty, "file", "f", prettyFile, "Name of the pretty file.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
