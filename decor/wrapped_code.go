package decor

// Function `wrappedTargetCode` is only used to show the structure of
// the target code after being decorated by the decorator.
// It is not actually executed.
//
// If the decorated target code causes a panic at runtime, Go's runtime will
// print the exception stack. The exception stack usually shows the path and line number
// of the relevant file where the call was made. This path is not friendly to
// decorator's real wrapped code.
// decorator solves this problem by pointing the wrapped function's associated code to
// this file at compile time.
//
// You shouldn't pay too much attention to the presence of this file error in the exception stack,
// but rather focus on the original target code.

func wrappedTargetCode( /* in1, in2, ... */ ) /* (out1, out2, ...) */ {
	varDecorContext := Context{
		Kind:       KFunc, // KFunc / KMethod
		TargetName: "",    // wrapped function/method name
		Receiver:   nil,   // wrapped method receiver
		TargetIn:   []any{ /*in1, in2, ....*/ },
		TargetOut:  []any{ /*out1, out2, ....*/ },
	}
	varDecorContext.Func = func() {
		/* varDecorContext.TargetOut[0], varDecorContext.TargetOut[1], ... = */ func( /* in1, in2, ... */ ) /* (out1, out2, ...) */ {

			// Here is code for the original target.

		}( /* varDecorContext.TargetIn[0], varDecorContext.TargetIn[1], ... */ )
	}

	/* decoratorFunc(varDecorContext, ...#{key: value}) */ // execute the decorator function

	return /* varDecorContext.TargetOut[0], varDecorContext.TargetOut[1], ... */
}
