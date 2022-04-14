# Differences between custom loggers

```
If zap is used normally, we can only print it on the console
But if we want the log to print to the file, we need to customize it
```

# zapcore.NewCore
```go
func NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core {
	return &ioCore{
		LevelEnabler: enab,
		enc:          enc,
		out:          ws,
	}
}
```
We can see that there are three parameters passed. One is encoder, writesyncer and levelenabler. The return value is core

#Zapcore log level

```go
const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)
```
We can see that it is an enumeration type to save

#Log cutting

Because our zap itself does not provide log cutting, but if we have many online logs, we need relevant cutting

We can use it at this time
```go
func getLogWriter1() zapcore.WriteSyncer {
    lumberjackLogger := &lumberjack.Logger{
        Filename: "./test.log",
        MaxSize: 10,
        MaxAge: 30,//Backup days

MaxBackups: 5,
        Compress: false,
    }
    return zapcore.AddSync(lumberjackLogger)
}
```