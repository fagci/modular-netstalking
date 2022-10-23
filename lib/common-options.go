package lib

import (
	"flag"
	"os"
)

var (
	OutputCount  uint64
	GreppableOutput bool
    DictPath string
)

func init() {
    o, _ := os.Stdout.Stat()
    isCharDevice := (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
	flag.BoolVar(&GreppableOutput, "oG", isCharDevice, "output as json")
}
