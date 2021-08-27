// Package filemngt
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package filemngt

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	tildeChar = '~'
	dotChar   = '.'

	emptyString = ""
	tildeString = "~"
	dotString   = "."
)

type ExtractMode int

const (
	NativeEMode ExtractMode = iota
	UserEMode
	CustomEMode
	UndefinedEMode
)

const (
	NativeEModeTag    = "native"
	UserEModeTag      = "user"
	CustomEModeTag    = "custom"
	UndefinedEModeTag = "undefined"
)

func IsValid(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

func startWithTilde(s string) bool {
	if len(s) > 0 {
		return s[0] == tildeChar
	}

	return false
}

func IsTilde(s string) bool {
	if len(s) == 1 {
		return s == tildeString || s[0] == tildeChar
	}

	return false
}

func IsDot(s string) bool {
	if len(s) == 1 {
		return s == dotString || s[0] == dotChar
	}

	return false
}

func IsPathSeparator(s string) bool {
	if len(s) == 1 {
		return os.IsPathSeparator(s[0])
	}

	return false
}

func MakeDirIfNotExist(d string) error {
	var err error
	if _, err = os.Stat(d); os.IsNotExist(err) {
		err = os.Mkdir(d, os.ModeDir)
		if err == nil {
			return nil
		}
	}

	return err
}

func MakeDir(d string) error {
	err := os.Mkdir(d, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// Check that the existing path is a directory
		info, err := os.Stat(d)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return ErrExistsIsNotADirectory()
		}

		return nil
	}

	return err
}

// GetWorkingPWD get the working pwd
func GetWorkingPWD(f string) (string, error) {
	return getWorkingPWD(f, NativeEMode)
}
func getWorkingPWD(fn string, m ExtractMode) (string, error) {
	if len(fn) == 0 {
		return emptyString, nil
	}
	var bPath string
	var err error

	bPath = filepath.Dir(fn)

	if IsTilde(bPath) || IsTilde(fn) || startWithTilde(bPath) {
		if IsTilde(fn) {
			bPath = tildeString
		}
		bPath, err = Expand(bPath, m)
		if err != nil {
			return emptyString, ErrCanNotExpandPath(bPath, err.Error())
		}

	}
	if len(bPath) == 0 ||
		IsPathSeparator(bPath) ||
		IsDot(bPath) {
		bPath, err = os.Getwd()
		if err != nil {
			return emptyString, ErrCanNotGetPWD(bPath, err.Error())
		}
	}

	vPath := filepath.VolumeName(bPath)
	if len(vPath) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return emptyString, ErrCanNotGetPWD(pwd, err.Error())
		}
		bPath = filepath.Join(pwd, bPath)
	}

	return bPath, nil
}

func GetFilePath(f string) (string, error) {
	return getFilePath(f, NativeEMode)
}

func getFilePath(f string, m ExtractMode) (string, error) {
	var basePath string
	var fPath string
	var err error
	var needExpand = false

	if len(f) == 0 {
		return emptyString, nil
	}

	if startWithTilde(f) {
		needExpand = true
	}

	_, ffn := filepath.Split(f)

	basePath, err = getWorkingPWD(f, m)
	if err != nil {
		return emptyString, err
	}

	if needExpand && ffn == tildeString {
		ffn = emptyString
	}

	fPath = filepath.Join(basePath, ffn)

	//// ffn = fnName + fnExt
	//fnExt := filepath.Ext(ffn)
	//fn := strings.TrimSuffix(ffn, fnExt)
	//fmt.Println("basePath:", basePath)
	//fmt.Println("full-filename", ffn)
	//fmt.Println("filename-extension:", fnExt)
	//fmt.Println("filename", fn)
	//fmt.Println("filePath", fPath)
	//fmt.Println("------")

	return fPath, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func Expand(path string, mode ExtractMode) (string, error) {
	if len(path) == 0 {
		return emptyString, nil
	}

	if path[0] != tildeChar {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return emptyString, errors.New("cannot expand user-specific home dir")
	}

	dir, err := getHomeDir(mode)
	if err != nil {
		return emptyString, err
	}

	return filepath.Join(dir, path[1:]), nil
}

func getHomeDir(mode ExtractMode) (string, error) {
	switch mode {
	case NativeEMode:
		dir, err := os.UserHomeDir()
		if err != nil {
			return emptyString, err
		}

		return dir, nil
	case UserEMode:
		u, err := user.Current()
		if err != nil {
			return emptyString, err
		}

		return u.HomeDir, nil
	case CustomEMode:
		dir, err := CustomHomeDir()
		if err != nil {
			return emptyString, err
		}
		return dir, nil
	default:
		return emptyString, ErrModeNotDefined(getModeTag(UndefinedEMode))
	}
}

func getModeTag(m ExtractMode) string {
	return getModeList()[m]
}
func getModeList() []string {
	return []string{NativeEModeTag, UserEModeTag, CustomEModeTag, UndefinedEModeTag}
}
