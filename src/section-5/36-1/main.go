package main

import (
	"fmt"
	"strings"
)

func shout(ping, pong chan string) {
	for {
		s := <-ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}
}

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	// to keep the program going and give the program a chance to the goroutine to execute, we need a time.Sleep in the simplest approach!
	//time.Sleep(10 * time.Second)
	fmt.Println("Type something and press ENTER (enter Q too quit)")

	for {
		// print a prompt
		fmt.Print("->")

		// get user input
		var userInput string

		// read whatever the user types
		_, _ = fmt.Scanln(&userInput) // scan it into userInput

		/* convert it to lower case, soo it doesn't matter if they type Q or q: */
		if userInput == strings.ToLower("q") {
			break
		} else {
			ping <- userInput

			// wait for a response(wait for the pong channel)
			response := <-pong
			fmt.Println("Response: ")
		}
	}
}
