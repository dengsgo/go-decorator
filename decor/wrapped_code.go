package decor

func wrappedTargetCode( /* in1, in2, ... */) /* (out1, out2, ...) */ {
	varDecorContext := Context{
		Kind:       KFunc, // KFunc / KMethod
		Receiver:   nil,   // wrapped method receiver
		TargetName: "",    // wrapped function/method name
		TargetIn:   []any{ /*in1, in2, ....*/ },
		TargetOut:  []any{ /*out1, out2, ....*/ },
	}
	varDecorContext.Func = func() {
		/* varDecorContext.TargetOut[0], varDecorContext.TargetOut[1], ... = */ func( /* in1, in2, ... */) /* (out1, out2, ...) */ {

			// Here is code for the original target.

		}( /* varDecorContext.TargetIn[0], varDecorContext.TargetIn[1], ... */)
	}

	/* decoratorFunc(varDecorContext, ...#{key: value}) */ // execute the decorator function

	return /* varDecorContext.TargetOut[0], varDecorContext.TargetOut[1], ... */
}
