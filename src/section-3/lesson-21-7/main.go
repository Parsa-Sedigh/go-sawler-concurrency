package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

/*
pizzasFailed: Sometimes when we try to produce sth, we're gonna fail for whatever reason.
total is pizzas made and pizzas failed added together
*/
var pizzasMade, pizzasFailed, total int

/*
The only thing the Producer knows about, is the channel data which is going to receive orders for pizzas and the channel quit, which

tells us we're all done making pizzas, so quit. In other words, stop the pizzeria goroutine that's running in the background.
*/
type Producer struct {
	// the producer is gonna receive an order for pizza:
	data chan PizzaOrder

	/* when we're finished making pizzas, for whatever reason, we're gonna send some info to the quite channel and sth else is gonna receive
	it and do sth with it. But in this case, it's a channel of channels! */
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		// since it's possible to get 0 as the result of Intn() , we add 1 to it because we want to delay for at least 1 second:
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d\n", pizzaNumber)

		/* Let's assume that in most cases, we make the pizza successfully, but if we hit some arbitrary number, then we fail */
		rnd := rand.Intn(12) + 1 // a random number between 1 and 12
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizzas #%d. It will take %d seconds...\n", pizzaNumber, delay)

		// delay
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quite while making pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("*** Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we're making(keep track of number of current pizza that we're making)
	var i = 0

	// run forever or until we receive a quit notification. We're running for doing what? Trying to make pizzas!
	for {
		// try to make a pizza
		currentPizza := makePizza(i)

		/* Probably it's impossible for currentPizza to nil, but instructor is a suspicious kind of guy! */
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			// we tried to make a pizza(here we send sth to the `data` channel):
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				//close channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
		// decision
	}

}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch

	return <-ch
}

func main() {
	/* seed the random number generator. If we don't seed the random number generator, we'll get the same result every time.

	With time.Now().UnixNano(), we ensure that we don't get the same results everytime we run the program*/
	rand.Seed(time.Now().UnixNano())

	// print out a message saying program is starting
	color.Cyan("The pizzeria is open for business")
	color.Cyan("---------------------------------")

	// create a producer and we have to describe it using some kind of data structure
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	/* run the producer in the background which means run it as it's own goroutine(we need a function to run it in background)*/
	go pizzeria(pizzaJob)

	// create and run consumer(s)
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.Close()
			if err != nil {
				// we won't have error in this case, but let's check for it anyway:
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("-----------------.")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an ok day...")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day!")
	default:
		color.Green("It was a great day!")

	}
}
