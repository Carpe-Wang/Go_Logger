package main

import (
	"go.uber.org/zap"
	"net/http"
)

var logger *zap.Logger

func main() {
	initLogger()
	defer logger.Sync()
	SimpleHttpRequest("https://www.baidu.com")
}
func initLogger() {
	logger, _ = zap.NewProduction()
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
