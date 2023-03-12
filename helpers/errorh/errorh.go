package errorh

import (
	"os"

	"github.com/notefan-golang/exceptions"

	"github.com/sirupsen/logrus"
)

func getLogger() (logger *logrus.Logger, err error) {
	logger = logrus.New()
	outputFile, err := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger.SetOutput(outputFile)
	return
}

// Panic will panic if the given anyErr is not nil
func Panic(anyErr any) {
	if anyErr == nil {
		return
	}

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

func Log(err any) {
	if err != nil {
		logger, openFileError := getLogger()
		Panic(openFileError)
		logger.Error(err) // error will continue the program
	}
}

func LogFatal(err any) {
	if err != nil {
		logger, openFileError := getLogger()
		Panic(openFileError)
		logger.Fatal(err) // fatal will exit the program immediately
	}
}

func LogPanic(err any) {
	if err != nil {
		logger, openFileError := getLogger()
		Panic(openFileError)
		logger.Panic(err) // panic will panicing the program immediately
	}
}
