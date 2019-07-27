package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/alecthomas/kingpin"
	"github.com/michaelsauter/go-oas-server/commands"
)

var (
	app = kingpin.New(
		"go-oas-server",
		"go-oas-server - Generate Go server code from an OpenAPI 3 specification",
	).DefaultEnvars()

	versionCommand = app.Command(
		"version",
		"Show version.",
	)

	generateCommand = app.Command(
		"generate",
		"Generate server code.",
	)
	generateFileFlag = generateCommand.Flag(
		"file",
		"Specification file.",
	).Short('f').Default("openapi.json").String()
	generateOutputDirFlag = generateCommand.Flag(
		"output-dir",
		"Output directory.",
	).Short('o').Required().String()
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("Fatal error: %s - %s.", err, debug.Stack())
		}
	}()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch command {
	case versionCommand.FullCommand():
		fmt.Println("0.1.0")
	case generateCommand.FullCommand():
		err := commands.Generate(*generateFileFlag, *generateOutputDirFlag)
		if err != nil {
			log.Fatalf("Failed to generate: %s.", err)
		}
	}
}
