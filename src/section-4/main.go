package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/* The dining philosophers problem is is well known in computer science circles. Five philosophers, numbered from 0 to 4,
live in a house where the table is laid for them; each philosopher has their own place at the table.
Their only difficulty - besides those of philosophy - is that the dish served is a very different kind of spaghetti which has to be eaten
with two forks. There are two forks next to each plate, so that presents no difficulty. As a consequence, however, this means that
no two neighbors may be eating simultaneously.*/

// constants
// this indicates how many times is the philosopher going to get hungry and we assume it's gonna eat 3 times
const hunger = 3

// variables - philosophers
var philosophers = []string{"Plato", "Socrates", "Aristotle", "Pascal", "Locke"}
var wg sync.WaitGroup
var sleepTime = 1 * time.Second
var eatTime = 3 * time.Second
var thinkTime = 1 * time.Second

// since this variable is used by multiple goroutines running at the same time, we're gonna have race condition, so we need a mutex for it
var orderFinished []string
var orderMutex sync.Mutex

func dingingProblem(philosopher string, leftFork, rightFork *sync.Mutex) {
	defer wg.Done()

	// print a message
	fmt.Println(philosopher, "is seated.")

	time.Sleep(sleepTime)

	/* We need to lock both forks, but we need to do it once for each time this philosopher gets hungry */
	for i := hunger; i > 0; i-- {
		fmt.Println(philosopher, "is hungry")
		time.Sleep(sleepTime)

		/* lock both forks which will make this goroutine to stop at this point. It will block until it can get a lock for the forks. */
		leftFork.Lock()
		fmt.Printf("\t%s picked up the fork for to his left.\n", philosopher)
		rightFork.Lock()
		fmt.Printf("\t%s picked up the fork for to his right.\n", philosopher)

		// print a message
		fmt.Println(philosopher, "has both forks and is eating.")
		time.Sleep(eatTime)

		// Give the philosopher some time ot think
		fmt.Println(philosopher, "is thinking")
		time.Sleep(thinkTime)

		// unlock the mutexes
		rightFork.Unlock()
		fmt.Printf("\t%s put down the fork on his right.\n", philosopher)
		leftFork.Unlock()
		fmt.Printf("\t%s put down the fork on his left.\n", philosopher)

		time.Sleep(sleepTime)
	}

	// print out done message
	fmt.Println(philosopher, "is satisfied.")
	time.Sleep(sleepTime)
	fmt.Println(philosopher, "left the table")

	// avoid race condition by using a mutex here:
	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher)
	orderMutex.Unlock()
}

func main() {
	// print intro
	fmt.Println("The Dining Philosophers Problem")
	fmt.Println("-------------------------------")

	wg.Add(len(philosophers))

	/* Make each philosopher try to eat. How are we going to do that? Spawn one goroutine for each philosopher and it will run until the philosophers
	finished eating.

	Everytime we go through this loop, the previous right fork(the one that's being created in the loop, that should actually be the left fork). So assign
	forkRight to forkLeft. But the first time through, we don't have a leftFork, so we the place to create that is outside the loop. So the very first
	mutex created, is the one created before we start spawning goroutines*/

	forkLeft := &sync.Mutex{}

	for i := 0; i < len(philosophers); i++ {
		// create a mutex for the right fork
		forkRight := &sync.Mutex{}
		// call(spawn) a goroutine
		go dingingProblem(philosophers[i], forkLeft, forkRight)

		/* Note: We never copy a mutex. Here these are both pointers, so we're not copying a mutex, we're making forkLeft equal to the pointer to an existing
		mutex. So it points to the same location in memory and doesn't copy it.*/
		forkLeft = forkRight
	}

	wg.Wait()

	fmt.Println("The table is empty")
	fmt.Println("=================")
	fmt.Printf("Order finsihed: %s\n", strings.Join(orderFinished, ", "))
}
