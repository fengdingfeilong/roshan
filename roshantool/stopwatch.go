package roshantool

import (
	"fmt"
	"time"
)

//Stopwatch time the function execute
//usage: defer Stopwatch("functionname")()
func Stopwatch(msg string) func() {
	start := time.Now()
	fmt.Printf("enter %s\n", msg)
	return func() {
		fmt.Printf("exit %s (%s)\n", msg, time.Since(start))
	}
}
