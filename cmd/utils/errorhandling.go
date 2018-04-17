package utils

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

////////////////////////// Exported funcs //////////////////////////

func RaiseCmdErr(cmd *cobra.Command, str string) error {
	return fmt.Errorf("%s.\nSee '%s -h' for help and examples", str, cmd.CommandPath())
}

func RaiseErr(str string) error {
	return fmt.Errorf("ERROR >> %s", str)
}

func CheckErr(str string) {
	fmt.Println(RaiseErr(str))
	os.Exit(1)
}
