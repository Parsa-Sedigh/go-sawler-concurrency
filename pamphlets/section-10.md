# Section 10: 10. Testing 

## 82-1. What we'll cover in this section
- update our models to be more testable
- test routes
- test the renderer
- test handlers

## 83-2. Setting up our tests
Let's set up a testing environment. When we're testing a web app, one of the things we have to do is to duplicate the environment that the 
various parts of our application run in and that's particularly true for handlers but it's true for most parts of our application.

The setup_test.go should be named exactly with this particular name and that file will run before our actual tests run and it sets up our environment.

TestMain is a special function used by go and it will run before our actual test run. In fact it will run our tests for us.

Note: We don't wanna use redis for a unit test. Because redis has been fully tested and we don't have to worry about that at all.

## 84-3. Testing Routes
Let's write a test to make sure that all of the routes in the project are actually in the project(the only thing we care about is 
we do in fact have these routes registered?)!

For this write routes_test.go .

To run the test, from the root of the project, run:
```shell
cd cmd/web
go test -v .
```

## 85-4. Testing the Renderer
Create render_test.go

## 86-5. Modifying the data package to make it testable
We want to make the data package more testable because the way it is right now, we need to have a database running to run unit tests nad that's
not good.

Currently if we want to run our unit tests, then we have to have a DB running. We might spin up a docker image with a known version of the database.
This is not good.

To fix this, create a file named interfaces in data directory.

By writing the UserInterface and PlanInterface and using them, it becomes possible too modify the data package, so we can have things that satisfy
the UserInterface and the PlanInterface, **but don't actually talk to the DB**! and that will make things easier.

Now create `test-models.go`.

**If you want to be able to pass sth as nil in a parameter of a function, make that param a reference to the type. Because args that are value and not reference
can't be nil.**

Solomon is that in test-models.go , we stub the functions of User and Plan models that don't touch DB at all(we need to satisfy the UserInterface and PlanInterface).

Now we have a version of User model that do not use DB at all.

## 87-6. Implementing the PlanTest type

## 88-7. Getting started testing Handlers
Now that our models don't use DB at all, let's write tests for handlers.

Create handlers_test.go in cdm>web.

## 89-8. Testing the Login Handler
```shell
go test -coverprofile=coverarge.out

# this will fire up a web browser and show us what is actually happening in the tests. Anything that's' green is covered and red is not covered.
go tool cover -html=coverageout
```

Currently, if we run the handlers_test.go , in PostLoginPage, the err handling of PasswordMatches is green, so it means it got executed during test and
we had an error there. The problem is when we said to ourselves, our test models are completely divorced from the DB, they are entirely divorced from the DB!
If we call user.<...> and we got that user var by calling app.Models.User.GetByEmail() , that's(app.Models.User.GetByEmail()) calling my test function, but
when we call it directly: user.<sth like PasswordMatches>, the problem we run into is it's calling the actual production version of PasswordMatches(our real
source code) and of course it's never gonna find a valid password because that user var doesn't have any of it's fields populated(it's populated with
just dummy values), so the PasswordMatches will fail.

To fix this, use the second approach.

This means that anywhere in the handlers that we're calling a method on the user var **directly**, in other words, we're not calling app.Models.User.<method> ,
we have to change those.

But you might ask yourself: It's convenient to call methods on a variable of type User, but we can't at the moment. Is there another way to do this?
There is, the way you would do this is that you would not separate your test DB functionality and your production DB functionality the way we did. Instead,
you would use sth called the repository pattern.

## 90-9. Testing a handler that uses concurrency
Hard coded values in source code cause us grief when testing! To fix this, you can declare variables and overwrite the value of those variables in tests(for example
when setting up the test env like TestMain func in our case).
We did this in handlers.go by creating pathToManual and tmpPath variables.

Currently if we run TestConfig_SubscribeToPlan test, nothings gonna happen because in the func we're testing, we have 2 goroutines that don't have time to finish!
One solution could be write:
testApp.Wait.Wait()
but it doesn't work.

If you run the test, it will stop at that test function forever!

We have 2 problems:
1) app.SendEmail() which runs concurrently, never decrease the wait-group
2) we need to decrement the wait group in <-testApp.Mailer.MailerChan of setup_test.go

To make sure there are no race conditions when running tests(running `go test -race` always takes longer):

```shell
go test -race -v .
```