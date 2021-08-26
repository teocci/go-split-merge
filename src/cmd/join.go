// Package cmd
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/teocci/go-split-merge/src/core"
	"os"
)

const (
	joinName  = "join"
	joinShort = "join files into small files"
	joinLong  = `This command join large files to improve portability over email or github commits.`

	joinFSFlagName  = "source"
	joinFSFlagShort = "s"
	joinFSFlagDesc  = "it will be join into a number of smaller files."

	joinFDFlagName  = "destination"
	joinFDFlagShort = "d"
	joinFDFlagDesc  = "it will be join into a number of smaller files."
	joinFDDefault   = "."
)

var (
	// joinCmd represents the add command
	joinCmd = &cobra.Command{
		Use:     joinName,
		Short:   joinShort,
		Long:    joinLong,
		PreRunE: joinPreRun,
		RunE:    joinRunE,
	}

	joinFSN string
	joinFDN string
)

func init() {
	joinFDN = joinFDDefault

	app.AddCommand(joinCmd)
	joinCmd.Flags().StringVarP(&joinFSN, joinFSFlagName, joinFSFlagShort, joinFSN, joinFSFlagDesc)
	joinCmd.Flags().StringVarP(&joinFDN, joinFDFlagName, joinFDFlagShort, joinFDN, joinFDFlagDesc)
	_ = joinCmd.MarkFlagRequired(joinFSFlagName)
}

func joinPreRun(cmd *cobra.Command, args []string) error {
	if Version {
		fmt.Printf(versionTemplate, appName, version, commit)
	}

	if _, err := os.Stat(joinFSN); err != nil {
		_ = cmd.Help()
	}

	return nil
}

func joinRunE(cmd *cobra.Command, args []string) error {
	var err error
	//if _, err = os.Stat(joinFSN); err != nil {
	//	_ = cmd.Help()
	//	return err
	//}

	_, _, err = core.Join(joinFSN, joinFDN)
	if err != nil {
		return err
	}

	return nil
}
