package helper

import (
	"os"

	"github.com/sirupsen/logrus"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorLog(err error) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(openFileError)
		logger.SetOutput(file)
		logger.WithError(err).Error() // error will continue the program
	}
}

func ErrorLogFatal(err error) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorLogPanic(openFileError)
		logger.SetOutput(file)
		logger.WithError(err).Fatal() // fatal will exit the program immediately
	}
}

func ErrorLogPanic(err error) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(openFileError)
		logger.SetOutput(file)
		logger.WithError(err).Panic() // panic will panicing the program immediately
	}
}
