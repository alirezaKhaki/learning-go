package main

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

// A function that returns an error
func doSomething() error {
	return errors.New("something went wrong")
}

func main() {
	// Call the function and capture the error
	err := doSomething()

	if err != nil {
		// Print the error with the stack trace
		fmt.Printf("%+v\n", err)
	}
}
