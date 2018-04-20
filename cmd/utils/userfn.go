package utils

import (
	"bufio"
	"fmt"
	"os"
)

func PromptUser() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		CheckErr(err.Error())
	}
	fmt.Println()
}
