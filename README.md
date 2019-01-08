# check
A collection of atomic checks that can be applied to values.

Each of the check functions returns an error if the value doesn't pass the
check or nil if the check is passed.

The checks are typically very simple to the point where you might question
why not perform the check directly. The reason is that as functions they can
be composed and combined and then passed to other code to be called later.
They are used extensively for checking command line parameters.
