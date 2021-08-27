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
	splitName  = "split"
	splitShort = "Split files into small files"
	splitLong  = `This command split large files to improve portability over email or github commits.`

	splitFFlagName  = "filename"
	splitFFlagShort = "f"
	splitFFlagDesc  = "it will be split into a number of smaller files."

	splitFDFlagName  = "destination"
	splitFDFlagShort = "d"
	splitFDFlagDesc  = "directory where the file will be split into [default './split']"
	splitFDDefault   = ""

	splitSizeFlagName    = "size"
	splitSizeFlagShort   = "s"
	splitSizeFlagDesc    = "size in MB in each of the smaller files"
	splitSizeFlagDefault = 48

	errSizeLesserThanOneMB = "size < 1, size should be greater than 1 MB"
)

var (
	// splitCmd represents the add command
	splitCmd = &cobra.Command{
		Use:     splitName,
		Short:   splitShort,
		Long:    splitLong,
		PreRunE: splitPreRun,
		RunE:    splitRunE,
	}

	splitFN  string
	splitDDN string
	size     int
)

func init() {
	splitDDN = splitFDDefault
	size = splitSizeFlagDefault

	app.AddCommand(splitCmd)
	splitCmd.Flags().StringVarP(&splitFN, splitFFlagName, splitFFlagShort, splitFN, splitFFlagDesc)
	splitCmd.Flags().StringVarP(&splitDDN, splitFDFlagName, splitFDFlagShort, splitDDN, splitFDFlagDesc)
	splitCmd.Flags().IntVarP(&size, splitSizeFlagName, splitSizeFlagShort, size, splitSizeFlagDesc)

	_ = splitCmd.MarkFlagRequired(splitFFlagName)
}

func splitPreRun(cmd *cobra.Command, args []string) error {
	if Version {
		fmt.Printf(versionTemplate, appName, version, commit)
	}

	if _, err := os.Stat(splitFN); err != nil {
		_ = cmd.Help()
	}

	if size < 1 {
		_ = cmd.Help()
		return fmt.Errorf(errSizeLesserThanOneMB)
	}

	return nil
}

func splitRunE(cmd *cobra.Command, args []string) error {
	var err error

	err = core.Split(splitFN, splitDDN, size)
	if err != nil {
		return err
	}

	return nil
}
