# Section 5: 5. Channels, and another classic The Sleeping Barber problem

## 36-1. What we'll cover in this section
Channels are the preferred method of sharing memory. Go's approach to concurrency, is share memory by communicating, don't communicate
by sharing memory and this is achieved primarily through the use of channels.

Once you fire off a goroutine into the background as it were, you really have no way of directly communicating with it, apart from the use
of channels.

Once you open a channel, you must close itt, otherwise, you're gonna wind uup with a resource leak.

## 37-2. Introduction to channels
When you're done with channels, close them otherwise. otherwise you're gonna have a resource leak.

When you get: `fatal erroor: all goroutines are sleep - deadlock!`. It means  you're sending sth to a channel, but nothing is listening to that channel
to receive the sent value. There's no goroutines(all goroutines are sleep) to listen that channel.

## 38-3. The select statement
When we receive sth from a channel, we have a second param we can get from that channel and it tells us whether the receive value of the chan(first
return value) was a zero value because the chan is closed(which indicates this by being false).

If there's more than one case that the select can match, it chooses one at random! and there's a lot oof situations where that's useful.

The default case in select statement is useful for avoiding deadlocks. Which means if there's a situation where none of the channels in select statement
are listening, then the default case will stoop your program from crashing. 

## 39-4. Buffered Channels
Buffered channels are useful when you **know** how many goroutines you've launched, or, we want to limit the number of goroutines we launch, or
we want to limit the amount of work that's queued up.

In the code of this example, we're limiting the amount of work that's queued up.

The vast majority of times you're goon use unbuffered channels.

## 40-5. Getting started with the Sleeping Barber project


## 41-6. Defining some variables, the barber shop, and getting started with the code
We're gonna have the shop that's running in the background doing it's thing and then we're gonna have some barbers. We're gonna have each
barber running as it''s own goroutine and for doing this, we're gonna define a method on the BarberShop named addBarber.
## 42-7. Adding a Barber
To exit a goroutine, we can write `return`.


## 43-8. Starting the barbershop as a GoRoutine
## 44-9. Sending clients to the shop
When the main program exits, any existing goroutines, they just die.

## 45-10. Trying things out
Don't forget too run the program using `-race` flag!

We didn't use any mutexes and wait groups.