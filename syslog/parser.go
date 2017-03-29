package syslog

// Parses 1 line of syslog and returns a parsed map of values
func Parse(line string) map[string]string {
	data := make(map[string]string)
	data["raw_input"] = line
	return data
}
