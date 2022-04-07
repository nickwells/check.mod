package check

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// ValCk is the type of a check function
type ValCk[T any] func(v T) error

// ValOK always returns a nil error - it can be of use with a
// check.ByPos to allow any value in certain slice positions
func ValOK[T any](_ T) error {
	return nil
}

// ValEQ returns a function that will check that the value is
// equal to the limit
func ValEQ[T comparable](limit T) ValCk[T] {
	return func(v T) error {
		if v == limit {
			return nil
		}
		return fmt.Errorf("the value (%v) must equal %v", v, limit)
	}
}

// ValNE returns a function that will check that the value is
// not equal to the limit
func ValNE[T comparable](limit T) ValCk[T] {
	return func(v T) error {
		if v != limit {
			return nil
		}
		return fmt.Errorf("the value (%v) must not equal %v", v, limit)
	}
}

// ValGT returns a function that will check that the value is
// greater than the limit
func ValGT[T constraints.Ordered](limit T) ValCk[T] {
	return func(v T) error {
		if v > limit {
			return nil
		}
		return fmt.Errorf("the value (%v) must be greater than %v", v, limit)
	}
}

// ValGE returns a function that will check that the value is
// greater than or equal to the limit
func ValGE[T constraints.Ordered](limit T) ValCk[T] {
	return func(v T) error {
		if v >= limit {
			return nil
		}
		return fmt.Errorf("the value (%v) must be greater than or equal to %v",
			v, limit)
	}
}

// ValLT returns a function that will check that the value is less
// than the limit
func ValLT[T constraints.Ordered](limit T) ValCk[T] {
	return func(v T) error {
		if v < limit {
			return nil
		}
		return fmt.Errorf("the value (%v) must be less than %v", v, limit)
	}
}

// ValLE returns a function that will check that the value is less
// than or equal to the limit
func ValLE[T constraints.Ordered](limit T) ValCk[T] {
	return func(v T) error {
		if v <= limit {
			return nil
		}
		return fmt.Errorf("the value (%v) must be less than or equal to %v",
			v, limit)
	}
}

// ValBetween returns a function that will check that the value
// lies between the upper and lower limits (inclusive)
func ValBetween[T constraints.Ordered](low, high T) ValCk[T] {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to ValBetween:"+
			" the lower limit (%v) must be less than the upper limit (%v)",
			low, high))
	}

	return func(v T) error {
		if v < low {
			return fmt.Errorf(
				"the value (%v) must be between %v and %v - too small",
				v, low, high)
		}
		if v > high {
			return fmt.Errorf(
				"the value (%v) must be between %v and %v - too big",
				v, low, high)
		}
		return nil
	}
}

// ValDivides returns a function that will check that the value
// is a divisor of d
func ValDivides[T constraints.Integer](d T) ValCk[T] {
	return func(v T) error {
		if d%v == 0 {
			return nil
		}
		return fmt.Errorf("the value (%d) must be a divisor of %d", v, d)
	}
}

// ValIsAMultiple returns a function that will check that the value
// is a multiple of d
func ValIsAMultiple[T constraints.Integer](d T) ValCk[T] {
	return func(v T) error {
		if v%d == 0 {
			return nil
		}
		return fmt.Errorf("the value (%d) must be a multiple of %d", v, d)
	}
}
