// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//   - if there are no customers, the barber falls asleep in the chair
//   - a customer must wake the barber if he is asleep
//   - if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and
//     sits in an empty chair if it's available
//   - when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//     and falls asleep if there are none
//   - shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//     empty
//   - after the shop is closed and there are no clients left in the waiting area, the barber
//     goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.
package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

// variables(package-level variables)
var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond

// how long is the barber shop open?
var timeOpen = 10 * time.Second

func main() {
	/* seed our random number generator and we use this rando number generator with arrival rate, so the clients don't always arrive
	at the same interval*/
	rand.Seed(time.Now().UnixNano())

	// print welcome message
	color.Yellow("The sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels if we need any
	// a channel that we send clients to it:
	clientChan := make(chan string, seatingCapacity)

	// a chan that says everything is done and we can go home. So whenever we're done, we'll send a bool to this chan
	doneChan := make(chan bool)

	// create the barbershop data structure
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	// add barbers
	shop.addBarber("Frank")
	shop.addBarber("Gerard")
	shop.addBarber("Milton")
	shop.addBarber("Susan")
	shop.addBarber("Kelly")
	shop.addBarber("Pat")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	// start the barbershop as a goroutine
	go func() {
		/* We need to make sure that this goroutine stays open at least(because there may be still some clients after the timeOpen
		duration in the shop) timeOpen amount of time. So we block this function for duration that the shop is supposed to be open*/
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients
	i := 1

	go func() {
		for {
			// get a random number with average arrival rate
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.addClient(fmt.Sprintf("Client %d", i))
				i++
			}
		}
	}()

	// block until the barbershop is closed(block until we receive sth from the closed channel)
	<-closed
}
