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
	mergeName  = "merge"
	mergeShort = "merge files into the original file"
	mergeLong  = `This command merge all file parts into the original file`

	mergeFSFlagName  = "source"
	mergeFSFlagShort = "s"
	mergeFSFlagDesc  = "first file of the split file e.g. xxxx.pt00"

	mergeFDFlagName  = "destination"
	mergeFDFlagShort = "d"
	mergeFDFlagDesc = "File that will contain the merged file"
	mergeFDDefault  = "."
)

var (
	// mergeCmd represents the add command
	mergeCmd = &cobra.Command{
		Use:     mergeName,
		Short:   mergeShort,
		Long:    mergeLong,
		PreRunE: mergePreRun,
		RunE:    mergeRunE,
	}

	mergeSFN string
	mergeDFN string
)

func init() {
	mergeDFN = mergeFDDefault

	app.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&mergeSFN, mergeFSFlagName, mergeFSFlagShort, mergeSFN, mergeFSFlagDesc)
	mergeCmd.Flags().StringVarP(&mergeDFN, mergeFDFlagName, mergeFDFlagShort, mergeDFN, mergeFDFlagDesc)

	_ = mergeCmd.MarkFlagRequired(mergeFSFlagName)
}

func mergePreRun(cmd *cobra.Command, args []string) error {
	if Version {
		fmt.Printf(versionTemplate, appName, version, commit)
	}

	if _, err := os.Stat(mergeSFN); err != nil {
		_ = cmd.Help()
	}

	return nil
}

func mergeRunE(cmd *cobra.Command, args []string) error {
	var err error
	//if _, err = os.Stat(mergeSFN); err != nil {
	//	_ = cmd.Help()
	//	return err
	//}

	_, _, err = core.Merge(mergeSFN, mergeDFN)
	if err != nil {
		return err
	}

	return nil
}
