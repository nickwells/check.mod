package check

import (
	"fmt"
	"regexp"
	"strings"
)

// StringLength returns a function that will apply the supplied check func to
// the length of a supplied value and return an error if the check function
// returns an error
func StringLength[T ~string](cf ValCk[int]) ValCk[T] {
	return func(v T) error {
		lv := len(v)
		err := cf(lv)
		if err == nil {
			return nil
		}
		return fmt.Errorf("the length of the string (%d) is incorrect: %w",
			lv, err)
	}
}

// StringMatchesPattern returns a function that checks that the string
// matches the supplied regexp. The regexp description should be a
// description of the string that will match the regexp. The error returned
// will say that the string "should be: " followed by this description. So,
// for instance, if the regexp matches a string of numbers then the
// description could be 'numeric'.
func StringMatchesPattern[T ~string](re *regexp.Regexp, reDesc string) ValCk[T] {
	return func(v T) error {
		if !re.MatchString(string(v)) {
			return fmt.Errorf("%q should be: %s", v, reDesc)
		}
		return nil
	}
}

// StringHasPrefix returns a function that checks that the string has the
// supplied string as a prefix
func StringHasPrefix[T ~string](prefix string) ValCk[T] {
	return func(v T) error {
		if !strings.HasPrefix(string(v), prefix) {
			return fmt.Errorf("%q should have %q as a prefix", v, prefix)
		}
		return nil
	}
}

// StringHasSuffix returns a function that checks that the string has the
// supplied string as a suffix
func StringHasSuffix[T ~string](suffix string) ValCk[T] {
	return func(v T) error {
		if !strings.HasSuffix(string(v), suffix) {
			return fmt.Errorf("%q should have %q as a suffix", v, suffix)
		}
		return nil
	}
}
