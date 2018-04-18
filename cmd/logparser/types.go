package logparser

import "time"

type SyslogMsg struct {
	Priority  int
	Severity  int
	Facility  int
	Name      string
	Message   string
	Hostname  string
	Timestamp time.Time
}
