/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"slice2go/client"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slice2go",
	Short: "Generate go files from slice file",
	Long:  `根据slice文件生成go相关文件`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	//PersistentPreRun: func(cmd *cobra.Command, args []string) { },
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			clientFlag.File = file
			client.Do(clientFlag)
		}
	},
}

var (
	isDebug = false
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var clientFlag client.ClientFlag

func init() {
	logrus.SetLevel(logrus.InfoLevel)

	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.support.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Print debug messages.")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.support.yaml)")

	rootCmd.Flags().StringVarP(&clientFlag.Prefix, "output-dir", "", ".", "Create files in the directory DIR.")
	rootCmd.Flags().StringVarP(&clientFlag.Interface, "interface", "i", "", "Generate only the specified interface.")
	rootCmd.Flags().StringVarP(&clientFlag.Template, "template", "t", "", "Use the template file from the specified directory")
	rootCmd.Flags().BoolVarP(&clientFlag.ExcludeIce, "exclude-ice", "", false, "Not generate ice cpp/h file")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".slice2go")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
