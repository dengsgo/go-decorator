## v0.22.0 beta

- Fix: the issue of compilation errors caused by underline parameters. #10

## v0.21.0 beta

- Fix: Possible null pointers during interface conversion

## v0.20.0 beta

- Feat: Support comment `type T types` type declaration, decorator will automatically decorate proxy all methods with `T` or `*T` as the receiver.  
- Add: More detailed error message display.  
- Add: Optimize the code called by the decorator, now the error message will prompt the correct line number  
- Add: More test cases and usage examples  
- Fix: Possible null pointers  

## v0.15.0 beta

- Feat: `decor.Context` added fields `TargetName`、`Receiver`
- Add: `TargetName` the function or method name of the objective function.
- Add: `Receiver` the receiver of the objective function. If `ctx.Kind == decor.KFunc` (i.e. function type), with a value of nil

## v0.12.0 beta

- Add: `usages/methods` demonstration cases
- Fix: the issue of decorators with parameters reporting errors if negative numbers are used

## v0.11.0 beta

- Feat: support methods for using decorators
- Feat: Reconstructed the demonstration project
- Feat: `decor.Context` added 'KMethod' type
- Fix: the issue of `required` lint reporting errors if negative numbers are used

## v0.10.0 beta

- Feat: support for decorators with parameters  
- Feat: the `go:decor-lint` annotation is supported to constrain the behavior of the caller  
- Feat: two built-in types of parameter constraints are implemented: `required`、`nonzero`  
- Update: new Scanning Implementation for Decorator Functions

## v0.5.0 beta

- Feat: decorator support for decorating generic functions

## v0.3.0 beta

- Feat: add method Context.DoRef()  
- Feat: streamline Code  
- Feat: go mod tidy. No third-party dependencies now  
- Add: documentation changes  

## 0.2.0 beta

- Feat: check repeated decoration  
- Feat: check func type  
- Add: test case funIsDecorator  
- Add: test case importer  

## 0.1.0 beta

- Initial project  
- Basic available features  

