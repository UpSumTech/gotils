package main

import (
	"fmt"
	"github.com/sumanmukherjee03/gotestbed/syslog"
	"log"
	"os"
	"regexp"
)

const GO_VERSION = "go1.8"

func input() string {
	var line string = "Mar 29 11:41:22 suman-mbp com.apple.xpc.launchd[1] (com.apple.nowplayingtouchui): Service  only ran for 0 seconds. Pushing respawn out by 10 seconds."
	return line
}

func reverse(arr [5]int) [5]int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func init() {
	var (
		user   = os.Getenv("USER")
		gopath = os.Getenv("GOPATH")
		msg    string
	)
	if _, err := regexp.MatchString(".*go1.8.*", gopath); err != nil {
		msg = fmt.Sprintf("$GOPATH %s for $USER %s is pointing to a wrong version of go. You will need version %s", gopath, user, GO_VERSION)
		log.Fatal(msg)
	}
}

func main() {
	inArr := [5]int{1, 2, 3, 4, 5}
	outArr := reverse(inArr)
	fmt.Println(outArr)
	str := input()
	fmt.Println(syslog.Parse(str))
}
