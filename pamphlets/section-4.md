# Section 4: 4. A Classic Problem The Dining Philosophers

## 27-1. What we'll cover in this section
This problem is solved in many different ways.

Concurrency is a good solution to this problem. We'll solve this problem not with channels but with sync.waitGroup and sync.mutex .

## 28-2. Getting started with the problem
We need to pass some info to the goroutine() indicating what philosopher we're dealing with and we also need to pass 2 mutexes. One for the left fork
and one for the right fork.

Where do we create those mutexes?

One answer could be we create an array and put one new mutex for each fork, but that's gonna get complicated quickly.

Since  things are running concurrently, things don't occur in sequential order, in some case we have people picking  up  a fork ono on side and
having to wait for the other one to  be available.

## 29-3. Setting up our mutexes
## 30-4. Finishing up the code
## 31-5. Trying things out
Let's give philosophers time to think. To do this, we add a delay before philosopher
## 32-6. Adding a delay to let a philosopher think
## 33-7. Challenge Printing out the order in which the philosophers finish eating
## 34-8. Solution to challenge
## 35-9. Writing a test for our program