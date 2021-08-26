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

	splitFFlagName  = "joinFN"
	splitFFlagShort = "f"
	splitFFlagDesc  = "it will be split into a number of smaller files."

	splitSizeFlagName    = "size"
	splitSizeFlagShort   = "s"
	splitSizeFlagDesc    = "size in MB in each of the smaller files"
	splitSizeFlagDefault = 48

	makeTmpDirFlagName    = "tmp-dir"
	makeTmpDirFlagShort   = "t"
	makeTmpDirFlagDesc    = "make a temporal directory for the hash and smaller files"
	makeTmpDirFlagDefault = true

	errSizeLesserThanOneMB = "size < 1, size should be greater than 1 MB"
)

var (
	// splitCmd represents the add command
	splitCmd = &cobra.Command{
		Use:     splitName,
		Short:   splitShort,
		Long:    splitLong,
		PreRunE: splitPreRun,
		RunE:    splitRun,
	}

	filename   string
	size       int
	makeTmpDir bool
)

func init() {
	size = splitSizeFlagDefault
	makeTmpDir = makeTmpDirFlagDefault

	app.AddCommand(splitCmd)
	splitCmd.Flags().StringVarP(&filename, splitFFlagName, splitFFlagShort, filename, splitFFlagDesc)
	splitCmd.Flags().IntVarP(&size, splitSizeFlagName, splitSizeFlagShort, size, splitSizeFlagDesc)
	splitCmd.Flags().BoolVarP(&makeTmpDir, makeTmpDirFlagName, makeTmpDirFlagShort, makeTmpDir, makeTmpDirFlagDesc)

	_ = splitCmd.MarkFlagRequired(splitFFlagName)
}

func splitPreRun(cmd *cobra.Command, args []string) error {
	if Version {
		fmt.Printf(versionTemplate, appName, version, commit)
	}

	if size < 1 {
		_ = cmd.Help()
		return fmt.Errorf(errSizeLesserThanOneMB)
	}

	return nil
}

func splitRun(cmd *cobra.Command, args []string) error {
	var err error
	if _, err = os.Stat(filename); err != nil {
		return err
	}

	err = core.Split(filename, size, makeTmpDir)
	if err != nil {
		return err
	}

	return nil
}
