## example

这个目录示范了如何正确编写代码来使用 go-decorator 工具。  

| Project                           | Notes                                                     |
|-----------------------------------|-----------------------------------------------------------|
| [**single**](example/single)      | 这个一个单文件示例，装饰器定义和被装饰的函数都位于一个包内。这种情况无需考虑导入依赖包的问题，按示例代码使用即可。 | 
| [**packages**](example/packages)  | 该项目示例为装饰器定义和被装饰的函数不在同一个包内，需要使用匿名包导入。                      |
| [**datetime**](example/datetime)  | Guide 里演示示例所用到的完整代码                                       |
| [**emptyfunc**](example/emptyfunc) | 演示装饰器中调用和不调用`TargetDo()` 的区别                              |


