package roshantool

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

//PrintQuitStack print current stack trace when stopped.
//this function should only be used for debugging
func PrintQuitStack() {
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			stack := make([]byte, 8192)
			len := runtime.Stack(stack, true)
			s := fmt.Sprintf("QUIT\nCurrent stack\n%s", string(stack[:len]))
			fmt.Println(s)
			if InnerLog != nil {
				InnerLog("roshan: "+s, nil)
			}
		}
		os.Exit(0)
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
	signal.Notify(sigChan, syscall.SIGINT)
	signal.Notify(sigChan, syscall.SIGKILL)
}
