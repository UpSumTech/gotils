package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sumanmukherjee03/gotils/cmd/logparser"
)

var (
	rootShortDesc = "Gotils is a single utilty to run various mixes of commands"
	rootLongDesc  = `Gotils is a Flexible tool built with Go.
	It does a mix of various things to help with developer and operations related work.`
	cfgFile string
	dryrun  bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:              "gotils [sub]",
		Short:            rootShortDesc,
		Long:             rootLongDesc,
		TraverseChildren: true,
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotils.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&dryrun, "dryrun", "d", false, "dryrun")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("useDryrun", rootCmd.PersistentFlags().Lookup("dryrun"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	rootCmd.AddCommand(logparser.NewLogParser())
	rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".gotils")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
