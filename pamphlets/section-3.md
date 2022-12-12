# Section 3: 3. Race Conditions, Mutexes, and an Introduction to Channels

## 15-1. What we'll cover in this section
mutex = mutual exclusion which allows us to deal with and prevent race conditions, in other words it deals with shared resources and concurrent/parallel
goroutines. Here **shared resources means variables or some data structure(sth that can be accessed by at least 2 goroutines at the same time)** and
OFC if you have 2 things(or more) running in the background(goroutines) and they both try to hit the same bit of data, bad things can happen, you may have
unpredictable results and we deal with this by locking and unlocking the resource as necessary.

We can test for race conditions in go either when we run a program, just by adding a flag with the `go` command, or when we test a program, when we run
a unit test or an integration test or sth like that.

Race conditions happen when we have at least 2 goroutines. So it's never gonna happen when you have just your main program and 1 other goroutine, but when you have
at least 2 goroutines running concurrently and they try to access the same thing. You can actually run a program that has a race condition, it runs
exactly as you expected and you move on thinking everything is fine! But fortunately, to find race conditions, go lets us check for these either
when we run them, or when we run an actual test.

## 16-2. Race Conditions an example
If we run a program with:
```shell
go run -race .
```
If you run this on the lesson's code, it will give us: `WARNING: DATA RACE` and a data race takes place when you have concurrent goroutines that access
the same piece of data and because you're never sure which one's gonna finish first, you actually run into problems.

We can fix this using mutex and unless we use that -race flag when running the program, we probably have no idea that we have a race condition.

## 17-3. Adding sync.Mutex to our code
**Like wait group, you don't want to copy a mutex once it's been created.**

After adding mutex(which causes a thread-safe operation), run the program again: `go run -race .`

## 18-4. Testing for race conditions
You can run:
```shell
go test -race .
```
to run the main_test.go in this lesson's code,  since we have 2 goroutines running in the background and both accessing the same data(msg variable),
we would get `WARNING: DATA RACE`. Remember when you run the test without -race flag, it passes, but with that flag, it shows there's problem.
So you don't have to test for race condition by **running** your program, you actually write tests and add that -race flag to make sure everything is gonna
behave as expected

## 19-5. A more complex example
## 20-6. Writing a test for our weekly income project
## 21-7. ProducerConsumer - Using Channels for the first time
## 22-8. Getting started with the Producer - the pizzeria function
## 23-9. Making a pizza the makePizza function
## 24-10. Finishing up the Producer code
## 25-11. Creating and running the consumer ordering a pizza
## 26-12. Finishing up our ProducerConsumer project