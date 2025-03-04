package check

// Aggregator is the type of an interface offering an Aggregate function and
// a Test function. The expectation is that the Aggregate function will be
// called over a range of values and will produce some kind of aggregate
// value which will be tested later with the test function.
type Aggregator[T any] interface {
	Aggregate(val T) error
	Test() error
}

// Counter implements the Aggregator interface. It counts the values that
// pass the value test and applies the count test
type Counter[T any] struct {
	valueTest ValCk[T]
	countTest ValCk[int]
	count     int
}

// NewCounter returns an instance of a Counter with the valueTest and
// countTest functions set from the parameters. If either test is nil a panic
// is generated.
func NewCounter[T any](valueTest ValCk[T], countTest ValCk[int]) *Counter[T] {
	if valueTest == nil {
		panic("no value test function has been given")
	}

	if countTest == nil {
		panic("no count test function has been given")
	}

	return &Counter[T]{
		valueTest: valueTest,
		countTest: countTest,
	}
}

// Aggregate counts the number of values that pass the valueTest
func (c *Counter[T]) Aggregate(val T) error {
	if c.valueTest(val) == nil {
		c.count++
	}

	return nil
}

// Test applies the testFunc to the count
func (c Counter[T]) Test() error {
	return c.countTest(c.count)
}
