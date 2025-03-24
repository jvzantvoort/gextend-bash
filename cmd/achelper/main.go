// bash autocomplete helper
package main

import (
	"fmt"
	"os"
)

func main() {
	firstArgument := os.Args[1]
	if firstArgument == "ssh" {
		err := SecureShellHelper()
		if err != nil {
			fmt.Println(err)
		}
	}
}
