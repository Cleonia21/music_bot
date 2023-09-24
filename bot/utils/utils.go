package utils

import (
	"fmt"
	"github.com/withmandala/go-log"
	"os"
)

func NewLogger(pwd string) (logger *log.Logger, err error) {
	file, err := os.Create(pwd)
	if err != nil {
		fmt.Println("logger " + pwd + " not created")
		return
	}
	logger = log.New(file)
	logger.WithDebug()
	//logger.WithColor()
	return
}
