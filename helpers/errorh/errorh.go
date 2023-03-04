package errorh

import (
	"os"

	"github.com/notefan-golang/exceptions"

	"github.com/sirupsen/logrus"
)

func Panic(anyErr any) {
	if anyErr != nil {
		err, ok := anyErr.(error)

		if ok {
			httpErr, ok := err.(*exceptions.HTTPError)
			if ok {
				panic(httpErr)
			}
		} else {
			panic(err)
		}
	}
}

func Log(err any) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		Panic(openFileError)
		logger.SetOutput(file)
		logger.Error(err) // error will continue the program
	}
}

func LogFatal(err any) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		Panic(openFileError)
		logger.SetOutput(file)
		logger.Fatal(err) // fatal will exit the program immediately
	}
}

func LogPanic(err any) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		Panic(openFileError)
		logger.SetOutput(file)
		logger.Panic(err) // panic will panicing the program immediately
	}
}
