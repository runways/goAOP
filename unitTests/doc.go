// Package unitTests. This package only used for AOP unit test
//
// In this package, there has a struct named FirstStruct. It has
// an unExport variable, named `name`, also has an export variable
// named `Age`. When AOP executes, these variables should be output normally.
//
// And It has two function, one has a comment used AOP, another has not AOP
// comment.
//
// '@middleware-a' is a test AOP comment. A valid AOP comment specification is that,
// starts with '@', and connect with a short word. e.g. '@middleware-a' or '@trace'.
package unitTests
