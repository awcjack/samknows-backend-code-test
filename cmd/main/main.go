package main

import (
	"log"
	"os"

	"github.com/awcjack/samknows-backend-code-test/app"
	"github.com/awcjack/samknows-backend-code-test/infrastructure/reader"
	"github.com/awcjack/samknows-backend-code-test/infrastructure/writer"
	"github.com/urfave/cli/v2"
)

func main() {
	// declare io reader and io writer that access filesystem
	ioReader := reader.NewIOReader()
	ioWriter := writer.NewIOWriter()
	// declare application that use io reader and io writer
	app := app.NewApplication(ioReader, ioWriter)

	cliApp := &cli.App{
		Name:      "performance-analyser",
		Usage:     "application that analyse the download performance and find the under-performing period",
		UsageText: "performance-analyser",
		Action: func(*cli.Context) error {
			err := app.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
