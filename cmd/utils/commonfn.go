package utils

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

////////////////////// Exported fns /////////////////////

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			CheckErr(err.Error())
		}
		viper.SetConfigType("yaml")
		viper.SetConfigFile(filepath.Join(home, ".gotils.yml"))
	}
	if err := viper.ReadInConfig(); err != nil {
		CheckErr(err.Error())
	}
}
