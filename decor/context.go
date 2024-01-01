package decor

// This file defines the context required for the decorator.
//
// If the function defined is of type func (* decor. Context), it is a decorator function,
// Can be used to decorate any top-level function.
//
// On top-level functions, use decorator functions through go single line annotations.
// For example:
//
// ```go
//   //go:decor decorHandlerFunc
//   func myFunc1() {
//	   log.Println("call myFunc1")
//   }
// ```
//
// The function myFunc1 declares the use of the decorator decorHandlerFunc The `go-decorator` tool
// will rewrite the target to inject decorHandlerFunc code during compilation.
// All of this is automatically completed at compile time!

// TKind is target types above and below the decorator
type TKind uint8

const (
	KFunc   TKind = iota // top-level function
	KMethod              // method
)

// Context The context of the decorator.
//
// The input and output parameters of the target function and the execution of
// the target method can be obtained through this context.
//
// Use TargetDo() to call the target function.
// If TargetDo() is not called in the decorator function, it means that the target
// function will not be called, even if you call the decorated target function in your code!
// At this point, the objective function returns zero values.
//
// Before TargetDo(), you can modify TargetIn to change the input parameter values.
// After TargetDo(), you can modify TargetOut to change the return value.
//
// You can only change the value of the input and output parameters. Don't try to change
// their type and quantity, as this will trigger runtime panic!!!
type Context struct {
	// Target types above and below the decorator
	Kind TKind

	// The input parameters of the decorated function
	TargetIn,

	// TargetOut : The result parameters of the decorated function
	TargetOut []any

	// The function or method name of the target
	TargetName string

	// If Kind is 'KMethod', it is the Receiver of the target
	Receiver any

	// The Non-parameter Packaging of the Objective Function // inner
	Func func()

	// The number of times the objective function was called
	doRef int64
}

// TargetDo : Call the target function.
//
// Calling this method once will automatically increment doRef by 1.
//
// Any problem can trigger panic, and a good habit is to capture it
// in the decorator function.
func (d *Context) TargetDo() {
	d.doRef++
	d.Func()
}

// DoRef gets the number of times an anonymous wrapper class has been executed.
// Usually, it shows the number of times TargetDo() was called in the decorator function.
func (d *Context) DoRef() int64 {
	return d.doRef
}
