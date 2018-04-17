package logparser

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

var (
	parseShortDesc = "Parses logs from a file"
	parseLongDesc  = `Parses logs from a file for known types.
		Current known types are : [syslog]`
	parseExample = `
	### Available commands for parse
	gotils parse syslog`
	supporttedParsers = map[string]func(line string) string{
		"syslog": parseSyslog,
	}
	src string
)

func NewLogParser() *cobra.Command {
	validLogTypes := []string{
		"syslog",
	}

	cmd := &cobra.Command{
		Use:     "parse LOGTYPE",
		Short:   parseShortDesc,
		Long:    parseLongDesc,
		Example: parseExample,
		Args: func(cmd *cobra.Command, args []string) error {
			var found bool
			if len(args) == 0 {
				return utils.RaiseCmdErr(cmd, "Kind of log not valid")
			}
			if len(args) > 1 {
				return utils.RaiseCmdErr(cmd, "Too many args")
			}
			for _, v := range validLogTypes {
				if !found {
					found = v == args[0]
				}
			}
			if !found {
				return utils.RaiseCmdErr(cmd, "Wrong type of log parser provided")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			readLogFile(args[0], src)
		},
	}

	cmd.Flags().StringVarP(&src, "src", "s", "", "Full path to the input file")
	cmd.MarkFlagRequired("src")
	return cmd
}

////////////////////////// Unexported funcs //////////////////////////

func readLogFile(kind string, fname string) {
	file, err := os.Open(fname)
	if err != nil {
		utils.CheckErr(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(parse(kind, scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		utils.CheckErr(err.Error())
	}
}

func parse(k string, s string) string {
	return supporttedParsers[k](s)
}
