package check_test

import (
	"fmt"

	"github.com/nickwells/check.mod/check"
)

// Example demonstrates how check functions might be used. It sets up two
// collections of checks on a slice of strings, the first collection should
// all pass (return a nil error) and the second set should all fail
func Example() {
	s := []string{"hello", "world"}

	passingChecks := []check.StringSlice{
		check.StringSliceLenEQ(2),
		check.StringSliceContains(
			check.StringEquals("hello"),
			"the list of strings must contain 'hello'"),
	}

	for _, c := range passingChecks {
		if err := c(s); err != nil {
			fmt.Println("unexpected error: ", err)
			return
		}
	}
	fmt.Println("All checks expected to pass, passed")

	failingChecks := []check.StringSlice{
		check.StringSliceLenEQ(3),
		check.StringSliceNot(
			check.StringSliceNoDups,
			"the list of strings must contain duplicates"),
	}

	var someCheckPassed bool
	for i, c := range failingChecks {
		if err := c(s); err == nil {
			fmt.Println("unexpected check success: ", i)
			someCheckPassed = true
		}
	}
	if !someCheckPassed {
		fmt.Println("All checks expected to fail, failed")
	}

	// Output:
	// All checks expected to pass, passed
	// All checks expected to fail, failed
}
