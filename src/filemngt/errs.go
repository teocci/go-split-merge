// Package filemngt
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
package filemngt

import (
	"errors"
	"fmt"
)

const (
	errCanGetPWD    = "cannot find the pwd: %s"
	errNotValidFile = "%s is not valid"
	errFileCannotBeOpened    = "file cannot be opened: %s"
	errModeNotDefined        = "%s mode not defined"
	errEmptyEnvironment      = "HOMEDRIVE, HOMEPATH, or USERPROFILE are blank"
	errEmptyOutputForHomeDir = "blank output when reading home directory"
	errCanNotExpandPath      = "path %s cannot be expanded: %s"
	errExistsIsNotADirectory = "path exists but is not a directory"
)

func ErrFileCannotBeOpened(e string) error {
	return errors.New(fmt.Sprintf(errFileCannotBeOpened, e))
}

func ErrNotValidFile(f string) error {
	return errors.New(fmt.Sprintf(errNotValidFile, f))
}

func ErrCanGetPWD(e string) error {
	return errors.New(fmt.Sprintf(errCanGetPWD, e))
}

func ErrModeNotDefined(m string) error {
	return errors.New(fmt.Sprintf(errModeNotDefined, m))
}

func ErrEmptyEnvironment() error {
	return errors.New(errEmptyEnvironment)
}

func ErrEmptyOutputForHomeDir() error {
	return errors.New(errEmptyOutputForHomeDir)
}

func ErrExistsIsNotADirectory() error {
	return errors.New(errExistsIsNotADirectory)
}

func ErrCanNotExpandPath(p string, e string) error {
	return errors.New(fmt.Sprintf(errCanNotExpandPath, p, e))
}
