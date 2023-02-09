package helper

import (
	"os"

	"github.com/sirupsen/logrus"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogIfError(err error) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		PanicIfError(openFileError)
		logger.SetOutput(file)
		logger.WithError(err).Error() // error will continue the program
	}
}

func LogFatalIfError(err error) {
	if err != nil {
		logger := logrus.New()
		file, openFileError := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		PanicIfError(openFileError)
		logger.SetOutput(file)
		logger.WithError(err).Fatal() // fatal will exit the program immediately
	}
}
