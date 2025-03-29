package check

import (
	"fmt"
)

// SliceLength returns a function that will apply the supplied check func to
// the length of a supplied value and return an error if the check function
// returns an error
func SliceLength[S ~[]E, E any](cf ValCk[int]) ValCk[S] {
	return func(v S) error {
		lv := len(v)

		err := cf(lv)
		if err == nil {
			return nil
		}

		return fmt.Errorf("the length of the list (%d) is incorrect: %w",
			lv, err)
	}
}

// SliceAggregate returns a function that will apply the Aggregate method of
// the suplied Aggregator to the values in the slice and will then return the
// results of the Test func.
//
// Note that if any of the calls to Aggregate returns a non-nil error the
// aggregation will stop and the error will be returned without the Test
// function being called.
func SliceAggregate[S ~[]E, E any](a Aggregator[E]) ValCk[S] {
	return func(v S) error {
		for _, e := range v {
			if err := a.Aggregate(e); err != nil {
				return err
			}
		}

		return a.Test()
	}
}

// SliceAll returns a function that will apply the supplied check func
// to each of the elements of the slice in turn and if any one of them fails
// the test it's location (index) in the slice and the error will be returned
// as an error.
//
// It returns nil if all the entries pass the check.
func SliceAll[S ~[]E, E any](cf ValCk[E]) ValCk[S] {
	return func(v S) error {
		if len(v) == 0 {
			return nil
		}

		for i, e := range v {
			if err := cf(e); err != nil {
				return fmt.Errorf(
					"list entry: %d (%v) does not pass the test: %w",
					i, e, err)
			}
		}

		return nil
	}
}

// SliceAny returns a function that will apply the supplied check
// func to each of the elements of the slice in turn and if all of them fail
// the test it returns an error. The msg parameter should describe the check
// being performed. For instance, for a slice of strings, if the check is
// that the string length must be greater than 5 characters then the
// condition parameter should be:
//
//	"the string should be greater than 5 characters"
//
// It returns nil if any of the entries pass the supplied check.
func SliceAny[S ~[]E, E any](cf ValCk[E], msg string) ValCk[S] {
	return func(v S) error {
		for _, e := range v {
			if err := cf(e); err == nil {
				return nil
			}
		}

		return fmt.Errorf("no list entries pass the test: %s", msg)
	}
}

// SliceByPos returns a check function that checks that a given entry in the
// slice passes the corresponding supplied check func. If there are more
// entries in the slice than check functions then the final supplied check is
// applied to any remaining entries. If there are fewer entries than
// functions then the excess checks are not applied. Note that you could
// choose to combine this check with a ValLength by means of a ValCkAnd
// check.
func SliceByPos[S ~[]E, E any](cfs ...ValCk[E]) ValCk[S] {
	return func(v S) error {
		if len(cfs) == 0 {
			return nil
		}

		for i, e := range v {
			if i >= len(cfs) {
				i = len(cfs) - 1
			}

			if err := cfs[i](e); err != nil {
				return fmt.Errorf(
					"list entry: %d (%v) does not pass the test: %w",
					i, e, err)
			}
		}

		return nil
	}
}

// SliceHasNoDups checks that the list contains no duplicates
func SliceHasNoDups[S ~[]E, E comparable](v S) error {
	dupMap := make(map[E]int)
	for i, s := range v {
		if dup, ok := dupMap[s]; ok {
			return fmt.Errorf("duplicate list entries: %d and %d are both: %v",
				dup, i, s)
		}

		dupMap[s] = i
	}

	return nil
}
