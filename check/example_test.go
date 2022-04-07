package check_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
)

// Example demonstrates how check functions might be used. It sets up two
// collections of checks on a slice of strings, the first collection should
// all pass (return a nil error) and the second set should all fail. Note
// that the check functions called below each returns a function of type
// check.StringSlice. For instance check.StringSliceLenEQ(2) returns a
// check.StringSlice function that will check that the given slice is of
// length 2. This technique is used throughout the package.
func Example() {
	s := []string{"hello", "world"}

	passingChecks := []check.ValCk[[]string]{
		check.SliceLength[[]string, string](check.ValEQ(2)),
		check.SliceAny[[]string](
			check.ValEQ("hello"),
			"the list of strings must contain 'hello'"),
	}

	for _, c := range passingChecks {
		if err := c(s); err != nil {
			fmt.Println("unexpected error: ", err)
			return
		}
	}
	fmt.Println("All checks expected to pass, passed")

	failingChecks := []check.ValCk[[]string]{
		check.SliceLength[[]string, string](check.ValEQ(3)),
		check.Not(
			check.SliceHasNoDups[[]string],
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
