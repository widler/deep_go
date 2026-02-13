package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (me *MultiError) Append(errs ...error) {
	me.errs = append(me.errs, errs...)
}

func (e *MultiError) Error() string {
	var res string
	if len(e.errs) == 0 {
		return ""
	}
	res = fmt.Sprintf("%d errors occured:\n", len(e.errs))
	for _, item := range e.errs {
		res += fmt.Sprintf("\t* %s", item.Error())
	}
	return res + "\n"
}

func Append(err error, errs ...error) *MultiError {
	switch err := err.(type) {
	case *MultiError:
		err.Append(errs...)
		return err
	default:
		res := &MultiError{}
		if err != nil {
			res.Append(err)
		}
		res.Append(errs...)
		return res
	}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
