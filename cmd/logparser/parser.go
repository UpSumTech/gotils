package logparser

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

var (
	parseShortDesc = "Parses logs from a file"
	parseLongDesc  = `Parses logs from a file for known types.
		Current known types are : [syslog]`
	parseExample = `
	### Available commands for parse
	gotils parse syslog`
	supportedParsers = map[string]func(line string) SyslogMsg{
		"syslog": ParseSyslog,
	}
	src      string
	dbUser   string
	dbPasswd string
	dbName   string
	dbHost   string
	dbPort   string
	db       *sql.DB
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
		PreRun: func(cmd *cobra.Command, args []string) {
			dbUser = viper.GetString("logparser.db.user")
			dbPasswd = viper.GetString("logparser.db.passwd")
			dbName = viper.GetString("logparser.db.name")
			dbHost = viper.GetString("logparser.db.host")
			dbPort = viper.GetString("logparser.db.port")
		},
		Run: func(cmd *cobra.Command, args []string) {
			db = dbConn()
			defer db.Close()
			stmts := initDbStmts()
			for _, stmt := range stmts {
				dbStmtExec(stmt)
			}
			readLogFile(args[0], src)
		},
	}

	cmd.Flags().StringVarP(&src, "src", "s", "", "Full path to the input file")
	cmd.MarkFlagRequired("src")
	cmd.Flags().StringVarP(&dbUser, "dbuser", "", "", "User name to use with the database")
	cmd.Flags().StringVarP(&dbPasswd, "dbpasswd", "", "", "Password to use with the database")
	cmd.Flags().StringVarP(&dbName, "dbname", "", "", "Name of the database")
	cmd.Flags().StringVarP(&dbHost, "dbhost", "", "", "Host of the database")
	cmd.Flags().StringVarP(&dbPort, "dbport", "", "", "Port of the database")
	viper.BindPFlag("logparser.db.user", cmd.Flags().Lookup("dbuser"))
	viper.BindPFlag("logparser.db.passwd", cmd.Flags().Lookup("dbpasswd"))
	viper.BindPFlag("logparser.db.name", cmd.Flags().Lookup("dbname"))
	viper.BindPFlag("logparser.db.host", cmd.Flags().Lookup("dbname"))
	viper.BindPFlag("logparser.db.port", cmd.Flags().Lookup("dbname"))
	return cmd
}

func Parse(kind string, line string) SyslogMsg {
	return supportedParsers[kind](line)
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
		msg := Parse(kind, scanner.Text())
		sql := "INSERT INTO parsed_data SET priority=?,facility=?,severity=?,name=?,hostname=?,message=?,timestamp=?"
		dbStmtExecWithVals(sql, msg)
	}

	if err := scanner.Err(); err != nil {
		utils.CheckErr(err.Error())
	}
}

func dbConn() *sql.DB {
	db, err := sql.Open("mysql", dbUser+":"+dbPasswd+"@"+"tcp"+"("+dbHost+":"+dbPort+")"+"/"+dbName)
	if err != nil {
		utils.CheckErr(err.Error())
	}
	return db
}

func dbStmtExec(str string) {
	stmt, err := db.Prepare(str)
	if err != nil {
		utils.CheckErr(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		utils.CheckErr(err.Error())
	}
}

func dbStmtExecWithVals(str string, msg SyslogMsg) {
	stmt, err := db.Prepare(str)
	if err != nil {
		utils.CheckErr(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(msg.Priority,
		msg.Facility,
		msg.Severity,
		msg.Name,
		msg.Hostname,
		msg.Message,
		msg.Timestamp)
	if err != nil {
		utils.CheckErr(err.Error())
	}
}

func initDbStmts() []string {
	sqlStmts := []string{
		"DROP TABLE IF EXISTS parsed_data",
		fmt.Sprintf("CREATE TABLE `parsed_data` (`uid` INT(10) NOT NULL AUTO_INCREMENT, `priority` INT(10) NOT NULL, `facility` INT(10) NOT NULL, `severity` INT(10) NOT NULL, `name` VARCHAR(64) NOT NULL, `hostname` VARCHAR(64) NOT NULL, `message` VARCHAR(256) NOT NULL, `timestamp` DATETIME NOT NULL, PRIMARY KEY (`uid`))"),
	}
	return sqlStmts
}
