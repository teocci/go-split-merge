// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package core

import (
	"bufio"
	"fmt"
	"github.com/teocci/go-split-merge/src/filemngt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Join(src, dest string) (string, []string, error) {
	var parts []string
	var basePath, srcFFN, mergedFileName, mergedFilePath string
	var err error

	if filemngt.IsValid(src) {
		basePath, srcFFN = filepath.Split(src)
		srcExt := filepath.Ext(srcFFN)
		srcFN := strings.TrimSuffix(srcFFN, srcExt)

		if len(basePath) == 0 {
			basePath, err = os.Getwd()
			if err != nil {
				return mergedFileName, nil, filemngt.ErrCanGetPWD(err.Error())
			}
		}

		parent := filepath.Dir(basePath)
		destPath := filepath.Join(parent, dest)

		fmt.Println("basePath:", basePath)
		fmt.Println("srcFFN:", srcFFN)
		fmt.Println("srcExt:", srcExt)
		fmt.Println("srcFN:", srcFN)
		fmt.Println("parent-path:", parent)

		err = filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() || path == basePath {
				r, err := regexp.MatchString(regExFirstElement + srcFN, d.Name())
				if err == nil && r {
					parts = append(parts, d.Name())
					fmt.Println("f.name:", d.Name())
				}
			} else {
				return filepath.SkipDir
			}

			return nil
		})
		if err != nil {
			return mergedFileName, nil, err
		}

		mergedFileName = joinFilePrefix + srcFN
		mergedFilePath = filepath.Join(destPath, mergedFileName)
		_, err = os.Create(mergedFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// set the mergedFile to APPEND MODE!!
		// open files r and w
		mergedFile, err := os.OpenFile(mergedFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		// IMPORTANT! do not defer a mergedFile.Close when opening a mergedFile for APPEND mode!
		// defer mergedFile.Close()

		// Just information on which part of the new mergedFile we are appending
		var writePosition int64 = 0
		for i, part := range parts {
			partPath := filepath.Join(parent, part)
			partFile, err := os.Open(partPath)
			if err != nil {
				log.Fatal(err)
			}
			defer partFile.Close()

			partInfo, err := partFile.Stat()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Processing file:", partInfo.Name())

			// calculate the bytes size of each chunk
			// we are not going to rely on previous data and constant
			partSize := partInfo.Size()
			partBytes := make([]byte, partSize)

			//fmt.Println("Appending at position : [", writePosition, "] bytes")
			writePosition = writePosition + partSize

			// read into partBytes
			reader := bufio.NewReader(partFile)
			_, err = reader.Read(partBytes)
			if err != nil {
				log.Fatal(err)
			}

			// DON't USE ioutil.WriteFile, it will overwrite the previous bytes!
			// Instead, write/save buffer to disk
			// ioutil.WriteFile(mergedFileName, partBytes, os.ModeAppend)
			n, err := mergedFile.Write(partBytes)
			if err != nil {
				log.Fatal(err)
			}

			_ = mergedFile.Sync() //flush to disk

			// Free up the buffer for next cycle should not be a problem if the
			// part size is small, but can be resource hogging if the part size is huge.
			// Also, it is a good practice to clean up your own plate after eating
			partBytes = nil // reset or empty our buffer

			fmt.Println("Written ", n, " bytes")
			fmt.Println("Recombining part [", i, "] into : ", mergedFileName)
		}

		// Now, close the mergedFile
		mergedFile.Close()
	}

	return mergedFilePath, parts, nil
}