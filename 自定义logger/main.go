package main

import (
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	_ "net/http"
	"os"
)

var logger *zap.Logger

func main() {
	initLogger()
	defer logger.Sync()
	SimpleHttpRequest("https://www.baidu.com")
}
func initLogger() {
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.InfoLevel)
	logger = zap.New(core)
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log") //创建文件，用来存放日志内容
	//	file, _ := os.Create("/.test.log",, os.O_APPEND) //创建文件，用来存放日志内容,追加
	return zapcore.AddSync(file)
}

//func getLogWriter1() zapcore.WriteSyncer {
//	lumberjackLogger := &lumberjack.Logger{
//		Filename: "./test.log",
//		MaxSize: 10,
//		MaxAge: 30,
//		MaxBackups: 5,
//		Compress: false,
//	}
//	return zapcore.AddSync(lumberjackLogger)
//}
func getEncoder() zapcore.Encoder {

	encoderconfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, //设置日志时间
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderconfig) //设置我们需要的json格式
}
func SimpleHttpRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"Error fetching url.....",
			zap.String("url", url),
			zap.Error(err),
		)
	} else {
		logger.Info("success:",
			zap.String("code", resp.Status),
			zap.String("url", url),
		)
		resp.Body.Close()
	}
}
