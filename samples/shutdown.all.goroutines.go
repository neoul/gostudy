// This example shows how to shutdown multiple goroutines using
// a single channel. The code makes use of the fact that receiving
// on a closed channel does not block.
// https://www.godesignpatterns.com/2014/04/exiting-multiple-goroutines-simultaneously.html

package statements

import (
	"fmt"
	"time"
)

var (
	shutdown = make(chan struct{})
	done     = make(chan int)
)

func shutdownallgoroutines() {
	const n = 5

	// Start up the goroutines...
	for i := 0; i < n; i++ {
		i := i
		go func() {
			fmt.Println("routine", i, "waiting to exit...")
			select {
			case <-shutdown:
				done <- i
			}
		}()
	}

	time.Sleep(2 * time.Second)

	// Close the channel. All goroutines will immediately "unblock".
	close(shutdown)

	for i := 0; i < n; i++ {
		fmt.Println("routine", <-done, "has exited!")
	}
}
