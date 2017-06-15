package main

import (
	"github.com/Smartling/api-sdk-go"
	"github.com/kovetskiy/lorg"
)

func setLogger(client *smartling.Client, logger lorg.Logger, verbosity int) {
	switch verbosity {
	case 0:
		return

	case 1:
		client.SetInfoLogger(logger.Infof)

	case 2:
		client.SetDebugLogger(logger.Debugf)
	}
}
