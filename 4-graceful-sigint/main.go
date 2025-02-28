//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	done := make(chan bool, 1)

	// Create a process
	proc := MockProcess{}

	go func() {
		sig := <-sigs

		if sig == syscall.SIGINT {
			fmt.Println("try to gracefully stop")
			go func() {
				proc.Stop()
				done <- true
			}()

			go func() {
				sig := <-sigs

				if sig == syscall.SIGINT {
					done <- false
				}
			}()
		}

		if status := <-done; status {
			fmt.Println("gracefully stopped")
			os.Exit(0)
		} else {
			fmt.Println("process killed")
			os.Exit(1)
		}
	}()

	// Run the process (blocking)
	proc.Run()
}
