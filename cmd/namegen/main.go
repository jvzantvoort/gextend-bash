package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/jvzantvoort/gextend-bash/version"
)

func main() {
	version.HandleArgs(os.Args)

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()

	fmt.Println(name)
}
