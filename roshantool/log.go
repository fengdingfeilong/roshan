package roshantool

import (
	"log"
	"os"
)

var logger *log.Logger
var fp *os.File

//InnerLog record the inner log of roshan
var InnerLog func(info string, err error)

//CreateLog ...
func CreateLog(file string) {
	fp, _ = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	logger = log.New(fp, "", log.LstdFlags)
}

//CloseLog ...
func CloseLog() {
	if fp != nil {
		fp.Close()
	}
}

//PrintErr ...
func PrintErr(method, err string) {
	logger.Printf("method %s Error: %s", method, err)
}

//Print ...
func Print(v ...interface{}) {
	logger.Print(v)
}

//Println ...
func Println(v ...interface{}) {
	logger.Println(v)
}

//Printf ...
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v)
}
