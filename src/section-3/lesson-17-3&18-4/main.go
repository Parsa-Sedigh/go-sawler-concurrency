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

	go updateMessage("Hello universe!")
	go updateMessage("Hello cosmos!")

	wg.Wait()

	fmt.Println(msg)
}

//func updateMessage(s string, m *sync.Mutex) {
//	defer wg.Done()
//
//	/* By calling Lock, we have exclusive access to msg, nobody else can change it's value until we're done with it. */
//	m.Lock()
//	msg = s
//	m.Unlock()
//}
//
//func main() {
//	msg = "Hello world!"
//
//	var mutex sync.Mutex
//
//	wg.Add(2)
//
//	/* We're still not sure which one's gonna finish first, so we're not sure what the result would be after running the program because we haven't
//	actually waited for the first goroutine to finish before the other one does, so we might get Hello universe! or Hello cosmos!, but the important
//	thing here is that we're accessing data safely(there is no race condition). This is what's called a thread-safe operation.*/
//	go updateMessage("Hello universe!", &mutex)
//	go updateMessage("Hello cosmos!", &mutex)
//
//	wg.Wait()
//
//	fmt.Println(msg)
//}
