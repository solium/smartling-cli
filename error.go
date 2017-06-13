package main

import "fmt"

type Error struct {
	Cause       error
	Description string
	Args        []interface{}
}

func NewError(cause error, description string, args ...interface{}) Error {
	return Error{
		Cause:       cause,
		Description: description,
		Args:        args,
	}
}

func (err Error) Error() string {
	return fmt.Sprintf(
		"ERROR: %s\n\n%s",
		err.Cause,
		fmt.Sprintf(err.Description, err.Args...),
	)
}
