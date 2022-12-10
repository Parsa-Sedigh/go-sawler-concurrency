# Section 2: 2. Goroutines, the go keyword, and WaitGroups
## 8-1. What we'll cover in this section

## 9-2. Creating GoRoutines
Every go program even the most simple one, have at least 1 goroutine. The main function itself is a goroutine.

Goroutines are very lightweight threads, not the builtin hardware threads of a processor, but instead the ones specific to go itself.
They take very little memory. They're all managed as a group of goroutines that is called coroutines. They are all managed by the go scheduler and
it decides what runs when, how much processing time one gets, it takes care of all that magic for us in the background.

```go
package main

import "fmt"

func print(str string) {
	fmt.Println(str)
}

func main() {
	go print("1")
	
	print("2")
}
```
The result would be: 2 and 1 won't be printed out at all. Because the program executed so quickly that there was not sufficient time for that
goroutine(the one we spawned on first line of main function) to actually execute!

How do we fix this?

There's a couple of ways. 3 good ways and 1 bad way. The bad way is to use `time.Sleep(1 * time.Second)` after spawning the goroutine.

We can use wait groups which allows us to wait for certain things to happen and then continue once they've taken place.


## 10-3. WaitGroups to the rescue
The previous solution was a terrible solution because if we had more goroutines, eventually they would take more than 1 second and the second
print could be printed before some of those goroutines. We don't know how much we should wait.

It doesn't matter what order you spawn goroutines, the order of them being completed, can't be predicted. In other words,
if you have multiple goroutines running at the same time, even if they're running the same function, you have no guarantee as to what order the complete,
that's entirely decided by the go's scheduler.

Once you created a wait group, you shouldn't copy it and instead pass pointers to it to the functions that use that wait group. 

## 11-4. Writing tests with WaitGroups
## 12-5. Challenge working with WaitGroup
## 13-5.2 go-concurrency-0240-goroutines-4
## 14-6. Solution to Challenge