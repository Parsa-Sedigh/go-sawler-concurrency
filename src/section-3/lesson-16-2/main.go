package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {
	msg = "Hello world!"
	wg.Add(2)

	/* We have a race condition. These two goroutines, are both called and they're called in a particular order, but we have no idea which one's
	gonna finish first.
	So even though the second goroutine is spawned(called) after the first one, it could be finished before the first one. If we run the program again,
	the order of logs may change.

	These two calls are both running at the same time. They're concurrent and we have no idea which one's going to finish first and both
	of them, access the package-level variable msg and because we're not sure which call to updateMessage is going to finish first,
	we have no idea what the value of msg is going to be by the time our program terminates.*/
	go updateMessage("Hello universe!")
	go updateMessage("Hello cosmos!")

	wg.Wait()

	fmt.Println(msg)
}
