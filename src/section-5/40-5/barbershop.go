package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		/* Initially the barber is awake. He or she arrives at work, presumably they're awake when they arrive at work!*/
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		// this is what we do endlessly until the day is over:
		for {
			// if there are no clients, the barber goes to sleep
			if len(shop.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			/* We need to make sure that the shop is still open. Because the rules say: once the shop is closed, we don't accept any new clients.
			So you might think that we could use the `Open` field of BarberShop. We could do that with one barber, but the minute I have more  than one barber,
			we have a potential data race(race condition). Because we might have more than one goroutine trying to READ or WRITE to the open field of
			the BarberShop. So instead we can use the second return type when reading a channel and it's not the same data that is stored in BarberShop.
			*/

			client, shopOpen := <-shop.ClientsChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up", client, barber)
					isSleeping = false
				}

				// cut hair
				shop.cutHair(barber, client)
			} else {
				// shop is closed, so send the barber home and close this goroutine
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)

	/* Now since the barber is going home, that means the goroutine associated with that barber is either gone or just about to disappear, so
	that barber can't take any moore clients. */
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day.")
	close(shop.ClientsChan)
	shop.Open = false

	/* Wait(block) until all the barbers are done: */
	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	color.Green("--------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day and everyone has gone home.")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("*** %s arrives!", client)

	if shop.Open {
		select {
		case shop.ClientsChan <- client:
			color.Yellow("%s takes a seat in the waiting room", client)

		/* We couldn't send to shop.clientsChan because the buffered channel is full: */
		default:
			color.Red("The waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Red("The shop is already closed, so %s leaves", client)
	}
}
