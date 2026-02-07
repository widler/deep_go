package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	// need to implement
	errs []error
}

func (me *MultiError) Append(err error) {
	me.errs = append(me.errs, err)
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
		for _, item := range errs {
			err.Append(item)
		}
		return err
	default:
		res := &MultiError{}
		if err != nil {
			res.Append(err)
		}
		for _, item := range errs {
			res.Append(item)
		}
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
