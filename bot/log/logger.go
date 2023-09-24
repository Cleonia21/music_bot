package log

import (
	"github.com/withmandala/go-log"
	"os"
)

var Logger *log.Logger

func Init() {
	Logger = log.New(os.Stderr)
	Logger.WithColor()
	Logger.WithDebug()
}
