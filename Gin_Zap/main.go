package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	//r := gin.Default()
	r := gin.New()
	r.Use(GinLogger(logger), GinRecovery(logger, true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello Carpewang!")
	})
	r.Run()
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

//func getLogWriter() zapcore.WriteSyncer {
//	file, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//	return zapcore.AddSync(file)
//}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     // M
		MaxBackups: 5,     // 最大备份数量
		MaxAge:     30,    // 最大备份天数
		Compress:   false, // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}
func GinLogger(Logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		path := context.Request.URL.Path
		query := context.Request.URL.RawQuery
		context.Next()

		cost := time.Since(start)
		Logger.Info(path,
			zap.Int("status", context.Writer.Status()),
			zap.String("method", context.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", context.ClientIP()),
			zap.String("user-agent", context.Request.UserAgent()),
			zap.String("errors", context.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if _, ok := err.(*net.OpError); ok {
					if se, ok := err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(context.Request, false)
				if brokenPipe {
					logger.Error(context.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					context.Error(err.(error)) // nolint: errcheck
					context.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				context.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		context.Next()
	}
}
