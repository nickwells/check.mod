package check

import "fmt"

// MapLength returns a function that will apply the supplied check func to
// the length of a supplied value and return an error if the check function
// returns an error
func MapLength[M ~map[K]V, K comparable, V any](cf ValCk[int]) ValCk[M] {
	return func(v M) error {
		lv := len(v)
		err := cf(lv)
		if err == nil {
			return nil
		}
		return fmt.Errorf("the length of the map (%d) is incorrect: %w",
			lv, err)
	}
}

// MapKeyAggregate returns a function that will apply the Aggregate method of
// the suplied Aggregator to the keys in the map and will then return the
// results of the Test func.
//
// Note that if any of the calls to Aggregate returns a non-nil error the
// aggregation will stop and the error will be returned without the Test
// function being called.
func MapKeyAggregate[M ~map[K]V, K comparable, V any](a Aggregator[K]) ValCk[M] {
	return func(v M) error {
		for k := range v {
			if err := a.Aggregate(k); err != nil {
				return err
			}
		}

		return a.Test()
	}
}

// MapValAggregate returns a function that will apply the Aggregate method of
// the suplied Aggregator to the values in the map and will then return the
// results of the Test func.
//
// Note that if any of the calls to Aggregate returns a non-nil error the
// aggregation will stop and the error will be returned without the Test
// function being called.
func MapValAggregate[M ~map[K]V, K comparable, V any](a Aggregator[V]) ValCk[M] {
	return func(m M) error {
		for _, v := range m {
			if err := a.Aggregate(v); err != nil {
				return err
			}
		}

		return a.Test()
	}
}

// MapKeyAll returns a function that will apply the supplied check function
// to each key in the map and will return an error for the first key for
// which it fails.
//
// It returns nil if all the keys pass the supplied check
func MapKeyAll[M ~map[K]V, K comparable, V any](cf ValCk[K]) ValCk[M] {
	return func(m M) error {
		for k := range m {
			if err := cf(k); err != nil {
				return fmt.Errorf("map entry[%v], bad key: %w", k, err)
			}
		}
		return nil
	}
}

// MapValAll returns a function that will apply the supplied check function
// to each value in the map and will return an error for the first value for
// which it fails.
//
// It returns nil if all the values pass the supplied check
func MapValAll[M ~map[K]V, K comparable, V any](cf ValCk[V]) ValCk[M] {
	return func(m M) error {
		for k, v := range m {
			if err := cf(v); err != nil {
				return fmt.Errorf("map entry[%v], bad value: %w", k, err)
			}
		}
		return nil
	}
}

// MapKeyAny returns a function that will apply the supplied check function
// to each key in the map and will return an error if all of them fail the
// test. The msg parameter should describe the check being performed. For
// instance, for a map keyed with strings, if the check is that the string
// length must be greater than 5 characters then the condition parameter
// should be:
//
//	"the string should be greater than 5 characters"
//
// It returns nil if any of the keys pass the supplied check
func MapKeyAny[M ~map[K]V, K comparable, V any](cf ValCk[K], msg string) ValCk[M] {
	return func(m M) error {
		for k := range m {
			if err := cf(k); err == nil {
				return nil
			}
		}
		return fmt.Errorf("no map keys pass the test: %s", msg)
	}
}

// MapValAny returns a function that will apply the supplied check function
// to each value in the map and will return an error if all of them fail the
// test. The msg parameter should describe the check being performed. For
// instance, for a map keyed with strings, if the check is that the string
// length must be greater than 5 characters then the condition parameter
// should be:
//
//	"the string should be greater than 5 characters"
//
// It returns nil if any of the values pass the supplied check
func MapValAny[M ~map[K]V, K comparable, V any](cf ValCk[V], msg string) ValCk[M] {
	return func(m M) error {
		for _, v := range m {
			if err := cf(v); err == nil {
				return nil
			}
		}
		return fmt.Errorf("no map values pass the test: %s", msg)
	}
}
