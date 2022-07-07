/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package errors

import (
	"github.com/pkg/errors"
	"strconv"
)

type DefineError struct {
	code    int
	message string
}

func NewDefineError(code int, message string) DefineError {
	return DefineError{code: code, message: message}
}

func NewDefineErrorWithCode(code int) DefineError {
	return DefineError{code: code}
}

func (d DefineError) Error() string {
	return d.message
}

func (d DefineError) Code() int {
	return d.code
}

func (d DefineError) Message() string {
	return d.message
}

func (d DefineError) Equal(err error) bool {
	return d.EqualError(d, err)
}

// EqualError equal error
func (d DefineError) EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}

// String parse code string to error.
func String(e string) DefineError {
	if e == "" {
		return NewDefineErrorWithCode(0)
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return NewDefineErrorWithCode(500)
	}
	return NewDefineErrorWithCode(i)
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return NewDefineErrorWithCode(0)
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}
