// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/teocci/go-split-merge/src/filemngt"
	"github.com/teocci/go-split-merge/src/units"
)

const (
	tmpSplitDir        = "./split"
	hashFileNamePrefix = "hash-"
	hashFileExt        = ".json"
	partFileExtPrefix  = ".pt"
)

// Split method splits the files into part files of user defined lengths
func Split(filename string, dest string, size int) error {
	if filemngt.IsPathValid(filename) {
		var err error
		bufferSize := int64(units.KiB)           // 1 KB for optimal splitting
		partSize := int64(size * int(units.MiB)) // Size in MiB

		filePath, err := filemngt.FilePathE(filename)
		if err != nil {
			return err
		}
		fileStats, _ := os.Stat(filePath)
		basePath, fn := filepath.Split(filename)

		workPath, err := findWorkPath(basePath, dest)
		if err != nil {
			return err
		}
		fmt.Println("workPath:", workPath)

		pieces := int(math.Ceil(float64(fileStats.Size()) / float64(partSize)))
		nTimes := int(math.Ceil(float64(partSize) / float64(bufferSize)))

		file, err := os.Open(filePath)

		hashFileName := hashFileNamePrefix + fn + hashFileExt
		hashFilePath := filepath.Join(workPath, hashFileName)
		_, _ = os.Create(hashFilePath)
		hashFile, _ := os.OpenFile(hashFilePath, os.O_CREATE, 0644)

		var hashes Hashes
		for i := 0; i < pieces; i++ {
			partNumTag := fmt.Sprintf("%02d", i)
			partFileName := fn + partFileExtPrefix + partNumTag
			partFilePath := filepath.Join(workPath, partFileName)
			partFile, err := os.OpenFile(partFilePath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return filemngt.ErrCanNotOpenFile(partFilePath, err.Error())
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

			appendHash(&hashes, partFileName, partFileHash)

			partFile.Close()
		}

		appendHash(&hashes, fn, Hash(filename))
		saveHashData(hashFile, hashes)

		file.Close()
		hashFile.Close()
		fmt.Printf("Splitted successfully! Find the individual file hashes in %s", hashFileName)
	} else {
		return filemngt.ErrNotValidFile(filename)
	}

	return nil
}

func findWorkPath(path string, dest string) (string, error) {
	var destPath string
	var err error
	if len(dest) == 0 {
		fmt.Println("Make working dir:", tmpSplitDir)
		destPath = filepath.Join(path, tmpSplitDir)
		err = filemngt.MakeDirIfNotExist(destPath)
		if err != nil {
			return emptyString, filemngt.ErrCanNotMakeDir(destPath, err.Error())
		}
	} else {
		destPath, err = filemngt.FilePathE(path)
		if err != nil {
			return emptyString, filemngt.ErrCanNotExpandPath(path, err.Error())
		}
	}

	return destPath, nil
}

func appendHash(h *Hashes, name string, hash string) {
	hd := HashData{
		Name: name,
		Hash: hash,
	}
	h.Hashes = append(h.Hashes, hd)
}

func saveHashData(f *os.File, h Hashes) {
	hdJson, _ := json.MarshalIndent(h, "", " ")
	_, _ = f.WriteString(string(hdJson))
}
