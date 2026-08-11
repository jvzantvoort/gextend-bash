package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jvzantvoort/gextend-bash/utils"
)

func Die(instr ...string) {
	if len(instr) != 0 {
		fmt.Printf("%s\n", instr)
	}
	fmt.Printf("USAGE:\n\n")
	fmt.Printf("\t%s <exitcode> <message>\n\n", os.Args[0])
	os.Exit(1)
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		Die()
	}
	exitcode, err := strconv.Atoi(args[0])
	if err != nil {
		Die(fmt.Sprintf("Error %s is not a number\n", args[0]))
	}

	message := strings.Join(args[1:], " ")

	if exitcode == 0 {
		fmt.Print(utils.MakeStatus("SUCCESS", "%s", message))
	} else {
		fmt.Print(utils.MakeStatus("FAILURE", "%s", message))
	}
	fmt.Printf("\n")
	os.Exit(exitcode)
}
