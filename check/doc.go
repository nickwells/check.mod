/*

Package check provides a collection of atomic checks that can be applied to
values.

Each of the check functions returns an error if the value doesn't pass the
check or nil if the check is passed. For instance a function of type check.X
will typically take a variable of type X and return an error if the value
does not pass the check. Where a check is parameterised there is typically a
function which returns a check function as a closure.

The checks are typically very simple to the point where you might question
why not perform the check directly. The reason is that as functions they can
be composed and combined and then passed to other code to be called later.
They are used extensively for checking command line parameters.

Many of the types have a ...Not function that can be used to invert the
meaning of a check. Similarly, there are ...And and ...Or functions which can
be used to compose checks.

*/
package check
