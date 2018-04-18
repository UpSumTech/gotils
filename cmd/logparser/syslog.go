package logparser

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/jeromer/syslogparser/rfc5424"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

//////////////////////////// Exported fns /////////////////////////////

func ParseSyslog(line string) SyslogMsg {
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

//////////////////////////// Unexported fns /////////////////////////////

func parseRfc3164(line string) (SyslogMsg, error) {
	var msg SyslogMsg
	buff := []byte(line)
	p := rfc3164.NewParser(buff)
	err := p.Parse()
	if err != nil {
		return msg, utils.RaiseErr(err.Error())
	}
	data := p.Dump()
	msg = SyslogMsg{
		Priority:  data["priority"].(int),
		Severity:  data["severity"].(int),
		Facility:  data["facility"].(int),
		Name:      data["tag"].(string),
		Hostname:  data["hostname"].(string),
		Message:   data["content"].(string),
		Timestamp: data["timestamp"].(time.Time),
	}
	return msg, nil
}

func parseRfc5424(line string) (SyslogMsg, error) {
	var msg SyslogMsg
	buff := []byte(line)
	p := rfc5424.NewParser(buff)
	err := p.Parse()
	if err != nil {
		return msg, utils.RaiseErr(err.Error())
	}
	data := p.Dump()
	msg = SyslogMsg{
		Priority:  data["priority"].(int),
		Severity:  data["severity"].(int),
		Facility:  data["facility"].(int),
		Name:      data["app_name"].(string),
		Hostname:  data["hostname"].(string),
		Message:   data["message"].(string),
		Timestamp: data["timestamp"].(time.Time),
	}
	return msg, nil
}
