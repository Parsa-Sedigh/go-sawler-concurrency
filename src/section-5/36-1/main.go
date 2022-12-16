package main

import (
	"fmt"
	"strings"
)

/*
	ping is a receive only channel and pong is a send-only channel. Now specifying these instead oof just saying `chan string`, is not necessary,

butt it  prevents you from accidentally trying tto send too a receive-only channel and receive from send-only channel,
*/
func shout(ping <-chan string, pong chan<- string) {
	for {
		// listens to the ping channel:
		s := <-ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}
}

func main() {
	// create two channels:
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

		if strings.ToLower(userInput) == "q" {
			break
		} else {
			ping <- userInput

			// wait for a response(wait for the pong channel)
			response := <-pong
			fmt.Println("Response: ", response)
		}
	}

	// when we get here, it means user typed a q to get out or quit
	close(ping)
	close(pong)
}
