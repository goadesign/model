/*
Package expr contains the data structure and associated logic that are built
by the DSL. These data structures are leveraged by the code generation
package to produce the software architecture model and views.

The expression data structures implement interfaces defined by the Goa eval
package. The corresponding methods (EvalName, Prepare, Validate, Finalize
etc.) are invoked by the eval package when evaluating the DSL.
*/
package expr
