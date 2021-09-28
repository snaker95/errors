package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

/*
@Time : 2021/6/22 下午7:03
@Author : snaker95
@File : error
@Software: GoLand
*/

const (
	UnknownCode = -100
)

type Error struct {
	code    int    `json:"code"`
	message string `json:"message"`
	err     error  `json:"err"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	format := "[%d]%s - %s"
	args := []interface{}{
		e.code,
		e.message,
		e.err,
	}
	return fmt.Sprintf(format, args...)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Cause() error {
	return e.err
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); As(err, &se) {
		return se.code == e.code
	}
	return false
}

func New(code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

func WithMessage(err error, message string, code ...int) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); As(err, &se) {
		return withMessage(err, se.code, fmt.Sprintf("%s -> %s", message, se.message))
	}
	c := UnknownCode
	if len(code) > 0 {
		c = code[0]
	}
	return withMessage(err, c, message)
}
func withMessage(err error, code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
		err:     err,
	}
}

// New returns an error object for the code, message.
func Wrap(err error, code int, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		code:    code,
		message: message,
		err:     errors.Wrap(err, ""),
	}
}

// Code returns the http code for a error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 0
	}
	if se := FormatError(err); err != nil {
		return se.code
	}
	return UnknownCode
}

// FormatError try to convert an error to *Error.
// It supports wrapped errors.
func FormatError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); As(err, &se) {
		return se
	}
	return withMessage(err, UnknownCode, "FormatError")
}
