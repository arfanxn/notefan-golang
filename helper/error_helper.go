package helper

import (
	"notefan-golang/exceptions"
	"os"

	"github.com/sirupsen/logrus"
)

func ErrorPanic(anyErr any) {
	if anyErr != nil {
		err, ok := anyErr.(error)

		httpErr, ok := err.(*exceptions.HTTPError)
		if ok {
			panic(httpErr)
		}

		panic(err)
	}
}

func ErrorLog(err any) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(openFileError)
		logger.SetOutput(file)
		logger.Error(err) // error will continue the program
	}
}

func ErrorLogFatal(err any) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(openFileError)
		logger.SetOutput(file)
		logger.Fatal(err) // fatal will exit the program immediately
	}
}

func ErrorLogPanic(err any) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(openFileError)
		logger.SetOutput(file)
		logger.Panic(err) // panic will panicing the program immediately
	}
}
