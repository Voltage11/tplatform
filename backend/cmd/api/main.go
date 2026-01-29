package main

import "tplatform/pkg/logger"

func main() {
	// Logger
	appLogger := logger.New(logger.AppLoggerLevelInfo)


	appLogger.SetLevel(logger.AppLoggerLevelError)

}