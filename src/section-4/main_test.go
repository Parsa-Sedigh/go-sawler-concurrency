package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	//set the delay variables to 0 so the test will run fast(actually the reason we created these variables was to replace them in tests!)
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 0 * time.Second

	/* we want to make sure this test works many times, so let's run the program 100 times. If we don't use orderMutex, the test will PROBABLY fail. Why
	probably? Because the race condition shows up once in a while.*/
	for i := 0; i < 100; i++ {
		main()

		if len(orderFinished) != 5 {
			t.Error("wrong number of entries in slice")
		}

		// reset orderFinished
		orderFinished = []string{}
	}

}
