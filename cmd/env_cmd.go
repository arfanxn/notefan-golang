package cmd

import (
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/helpers/cmdh"
)

// GuessENV will guess env configuration based on the command line arguments
func GuessENV() {
	switch true {
	case cmdh.UserFirstArgIs("test"):
		config.LoadTestENV()
		break
	default:
		config.LoadENV()
	}
}
