# Section 5: 5. Channels, and another classic The Sleeping Barber problem

## 36-1. What we'll cover in this section
Channels are the preferred method of sharing memory. Go's approach to concurrency, is share memory by communicating, don't communicate
by sharing memory and this is achieved primarily through the use of channels.

Once you fire off a goroutine into the background as it were, you reAlly have no way of directly communicating with it, apart from the use
oof channels.

OOnce you open a channel, you must close itt, otherwise, you're gonna wind uup with a resource leak.

## 37-2. Introduction to channels
## 38-3. The select statement
## 39-4. Buffered Channels
## 40-5. Getting started with the Sleeping Barber project
## 41-6. Defining some variables, the barber shop, and getting started with the code
## 42-7. Adding a Barber
## 43-8. Starting the barbershop as a GoRoutine
## 44-9. Sending clients to the shop
## 45-10. Trying things out