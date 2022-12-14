package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string

	// how much do you make from this source? The amount is represented in this field:
	Amount int
}

func main() {
	// variable for bank balance:
	var bankBalance int // because we don't initialize the variable, it has a default zero value of 0 which is exactly what we want it to be!
	var balance sync.Mutex

	//print out starting values
	fmt.Printf("Initial account balance: %d.00", bankBalance)
	fmt.Println() // just print a blank line

	/* Define weekly revenue. How much money do you make in a week and where does it come from? and we'll have multiple sources
	of income, just to give us sth to work with.*/
	incomes := []Income{
		{Source: "Main job", Amount: 500},
		{Source: "Gifts", Amount: 10},
		{Source: "Part time job", Amount: 50},
		{Source: "Investments", Amount: 100},
	}

	wg.Add(len(incomes))

	// loop through 52 weeks which is 1 year and print out how much is made; also keep a running total
	for i, income := range incomes {
		// we inlined the function here, we could extract it
		go func(i int, income Income) {
			defer wg.Done()

			for week := 1; week <= 52; week++ {
				balance.Lock()
				// current bank balance:
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()

				fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, income.Amount, income.Source)
			}
		}(i, income)
	}

	wg.Wait()

	// print out final balance
	fmt.Printf("Final bank balance: $%d.00", bankBalance)
	fmt.Println()
}
