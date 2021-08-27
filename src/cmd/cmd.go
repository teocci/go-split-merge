// Package cmd
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const (
	appName  = "go-split-merge"
	appShort = "Split/Join files"
	appLong  = `This application split/join large files to improve portability over email or github commits.`
	version  = "v1.1"
	commit   = "0"
)

const (
	versionTemplate = "%s %s.%s\n"
)

var (
	// shaman provides the shaman cli/server functionality
	app = &cobra.Command{
		Use:           appName,
		Short:         appShort,
		Long:          appLong,
		PreRunE:       preRun,
		Run:           run,
		SilenceErrors: false,
		SilenceUsage:  true,
	}

	Verbose  = false // Run in verbose mode
	Version  = true // Print version info and exit
)

func run(ccmd *cobra.Command, args []string) {

}

func Execute() {
	err := app.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	//app.Flags().BoolVarP(&isSplit, "is-split", "s", isSplit, "If the zip file has been split")

	app.Flags().BoolVarP(&Verbose, "verbose", "v", Verbose, "Run in debug mode")
	app.Flags().BoolVarP(&Version, "version", "V", Version, "Print version info and exit")
}

func preRun(ccmd *cobra.Command, args []string) error {
	if Version {
		fmt.Printf(versionTemplate, appName, version, commit)
	}

	if len(args) == 0 {
		_ = ccmd.Help()
		os.Exit(0)
	}

	return nil
}
