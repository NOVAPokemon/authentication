package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	// RegisterHandler errors
	errorRegisteringConflictFormat = "tried to register %s again"

	// Login errors
	errorWrongPasswordFormat = "wrong password for user %s"

	// Refresh errors
	errorRefreshTooSoonFormat = "user %s tried to refresh token too soon"

	// Others
	errorHashPassword = "error hashing password"
)

// Register
func wrapRegisterHandlerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, RegisterName))
}

func newRegisterConflictError(username string) error {
	return errors.New(fmt.Sprintf(errorRegisteringConflictFormat, username))
}

// Login
func wrapLoginHandlerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, LoginName))
}

func newWrongPasswordError(username string) error {
	return errors.New(fmt.Sprintf(errorWrongPasswordFormat, username))
}

// Refresh
func wrapRefreshHandlerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, RefreshName))
}

func newRefreshTooSoonError(username string) error {
	return errors.New(fmt.Sprintf(errorRefreshTooSoonFormat, username))
}

// Others
func wrapHashPasswordError(err error) error {
	return errors.Wrap(err, errorHashPassword)
}
