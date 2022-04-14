#这个是学习gin和zap整合的笔记代码
查阅default源码我们可以发现
```go
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
```
它是采用Logger中间件。如果我们想要使用logger