package main

import (
	"io"
  "github.com/morriswinkler/uhrwerk/debug"
)

/*
type NfcConfig struct {
	Basedir string
}

type config struct {
	Nfc NfcConfig
}
*/

func Init(
  traceHandle io.Writer,
  infoHandle io.Writer,
  warningHandle io.Writer,
  errorHandle io.Writer) {

	debug.Init(traceHandle, infoHandle, warningHandle, errorHandle)
}

