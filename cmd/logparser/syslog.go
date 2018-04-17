package logparser

import (
	"fmt"

	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/jeromer/syslogparser/rfc5424"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func parseRfc3164(line string) (string, error) {
	buff := []byte(line)
	p := rfc3164.NewParser(buff)
	err := p.Parse()
	if err != nil {
		return line, utils.RaiseErr(err.Error())
	}
	return utils.ToJson(p.Dump()), nil
}

func parseRfc5424(line string) (string, error) {
	buff := []byte(line)
	p := rfc5424.NewParser(buff)
	err := p.Parse()
	if err != nil {
		return line, utils.RaiseErr(err.Error())
	}
	return utils.ToJson(p.Dump()), nil
}

func parseSyslog(line string) string {
	s, err := parseRfc3164(line)
	if err != nil {
		fmt.Println(err)
		s, err = parseRfc5424(line)
		if err != nil {
			fmt.Println(err)
			utils.CheckErr("Could not parse syslog line in Rfc3164 or Rfc5424 formats")
		}
	}
	return s
}
