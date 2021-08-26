// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package core

import (
	"fmt"
	"github.com/teocci/go-split-merge/src/filemngt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
)


// Split method splits the files into part files of user defined lengths
func Split(filename string, size int, makeTmpDir bool) error {
	if filemngt.IsValid(filename) {
		var err error
		bufferSize := int64(KiB)           // 1 KB for optimal splitting
		partSize := int64(size * int(MiB)) // Size in MiB

		filePath, err := filemngt.GetFilePath(filename)
		if err != nil {
			log.Fatal(err)
		}

		basePath, fn := filepath.Split(filename)

		fileStats, _ := os.Stat(filePath)

		pieces := int(math.Ceil(float64(fileStats.Size()) / float64(partSize)))
		nTimes := int(math.Ceil(float64(partSize) / float64(bufferSize)))

		file, err := os.Open(filePath)

		hashFileName := hashFileNamePrefix + fn
		var hashFilePath string
		var tmpDirPath string
		if makeTmpDir {
			fmt.Println("Make tmp dir:", tmpSplitDir)
			tmpDirPath = filepath.Join(basePath, tmpSplitDir)
			err := filemngt.MakeDirIfNotExist(tmpDirPath)
			if err != nil {
				log.Fatal(err)
			}

			hashFilePath = filepath.Join(tmpDirPath, hashFileName)
		}

		_, err = os.Create(hashFilePath)
		if err != nil {
			log.Fatal(err)
		}

		hashFile, err := os.OpenFile(hashFilePath, os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < pieces; i++ {
			partNumTag := fmt.Sprintf("%02d", i)
			partFileName := fn + hashFileExtPrefix + partNumTag
			partFilePath := filepath.Join(tmpDirPath, partFileName)
			partFile, err := os.OpenFile(partFilePath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Creating file:", partFileName)
			buffer := make([]byte, bufferSize)
			for j := 0; j < nTimes; j++ {
				_, rErr := file.Read(buffer)
				if rErr == io.EOF {
					break
				}

				_, wErr := partFile.Write(buffer)
				if wErr != nil {
					log.Fatal(wErr)
				}
			}

			partFileHash := Hash(partFilePath)
			s := partFileName + ": " + partFileHash + "\n"
			_, _ = hashFile.WriteString(s)
			partFile.Close()

		}
		s := "original-file-hash: " + Hash(filename) + "\n"
		_, _ = hashFile.WriteString(s)

		file.Close()
		hashFile.Close()
		fmt.Printf("Splitted successfully! Find the individual file hashes in %s", hashFileName)
	} else {
		return filemngt.ErrNotValidFile(filename)
	}

	return nil
}
