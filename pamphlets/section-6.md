# Section 6: 6. Final Project - Building a Subscription Service

## 46-1. What we'll cover in this section

## 47-2. Setting up a simple web application
In the main func, since we're sending mails, we need to use wait groups, because we're sending emails concurrently. Now if someone decides
to restart the application for whatever reason(maybe doing some maintenance or ...), you don't want to just say: Stop the application. Instead,
we want to wait until any mail that is waiting to be sent, gets delivered and then you stop the app and this is an ideal situation for wait group.

There are a couple of different drives for postgres, one is PQ but it says use the jackc's one instead which is named `pgconn`! So let's use that one and then
`pgx/v4`. For session management we use `scs/v2` . This package allows different stores for session data. You can store it in the database, in a cookie, but we 
want to use redis and for this, install: `scs/redisstore`. Also install the chi router: `go-chi/chi/v5`.

## 48-3. Setting up our Docker development environment
We're gonna use redis as a store for our session info and we're also gonna require some kind of mail server, a dummy mail server that will capture mail
for us and we're gonna use docker for that.

**Mailhog**(dummy mail server): it allows us to send email without actually sending email! Mailhog captures it and if offers a web interface where we can go and look 
at our email.

In docker-compose file, we specified a volume for storing the data for postgres in the `./db-data/postgresa/...`. Now in a lot of cases, 
docker will create that for you, but for safety, you can create a new folder at the root level of your project called db-data and `postgres` directory and
one for `redis`.

Then you can run:
```shell
docker compose up -d # -d runs this in the background
```

Now you want to open your favorite database client(for not having to use command line) like `beekeeper studio`.

## 49-4. Adding postgres
Now we need to setup the DB like connecting to it in code, but for this, we want to connect to it using the DB client and enter user and password and
the name of the DB you want to connect to(in dev, we're connecting to localhost as host).

Then connect to DB in code.

By saying:
```go
import (
    _ "import path"
)

```
We say even though this package isn't directly used in the code, we still need to have it there.

While the necessary containers are running, run the program. But we need the DSN. Now we could just set an environment variable and run the program,
but this is the point where we should start using `make`.

## 50-5. Setting up a Makefile
The process of creating a makefile on windows is different than on mac.

It's good practice when building a binary go app to set CGO_ENABLED to 0 if you're not using CGO.

DSN == data source name

Run make start in the same directory that the Makefile exist to run the app (like connect to the DB) 

## 51-6. Adding sessions & Redis
## 52-7. Setting up the application config
Create config.go in cmd>web directory.

## 53-8. Setting up a route & handler for the home page, and starting the web server
After finishing this lesson, we can start the server in background by saying: `make start` and stop by: `make stop`.

## 54-9. Setting up templates and building a render function
The `base.layout.gohtml` is the base layout for every template.

To render the templates, create a file named render.go . There, create `pathToTemplates` var and we did this so when we rite tests,
we don't run into the problem of having to somehow override a string that's a constant.

Now we need sth to store datta in that we're gonna pass to the templates for this, let's create `TemplateData` struct type which specifies the kind of
things we're gonna pass to templates(we might not use all of the fields in there).

We can use `interface{}` instead of `any` in go. It's functionally equivalent.

If we have one template that depends on a partial oro a layout oor anything like that, then I need to include all oof those templates
and direct the full path name to every one of them, when we go to parse the templates.

There are certain kinds of data we want to pass to every single template, so what we could we do to make sure that every template get that
data? `AddDefaultData` func.

## 55-10. Adding session middleware
We can't use our app until we set up some middleware that loads and saves the session on every req and we can do this with `scs` package.
In CMD/Web create middleware.go .

In order to use the middleware, in routes.go write: `mux.Use()` in `routes` func.

To start the web app, run: `make start` in the directory where the `Makefile` exists.

Now in the web browser, you can open `localhost`.

## 56-11. Setting up additional stub handlers and routes
After a user register we'll send him an activation email to verify that we have the valid email address and we'll have them activate their account
and that will be a GET req and the handler for this is ActivateAccount func.

Now run:
```shell
make restart
```
to restart the application.

## 57-12. Implementing graceful shutdown
Our app is gonna have a number of goroutines running iin the background. Some oof those goroutines are just gonna be listening to
channels and some are gonna be doing sth like sending gan email or generating invoice or sth like that and when you decide to stop this app, for
whatever reason, maybe you need to do some work on the server, maybe you need to implement a hotfix or ..., you need to stop the application
and if you just stop it, like running `make stop`, everything just stops, any running goroutines just die without finishing and that's not good
because you might not send an email that needs to go out or generate an invoice or ...

So it's a good practice to implement graceful shutdown. For this, create a func that will be running by a goroutine named `listenForShutdown` and it's
just running in the background and just listen for sth.

**SIGINT** means interrupt signal(int means interrupt).

The logs for this may not show up on windows because some you might be using a shell that don't get output in the terminal when they're running `make`
commands.

Now for example everytime we want to send an email in the background or generating an invoice or sth that's gonna be running in the background,
we'll increase the app's wait group and also we'll have a `defer app.Wait.Done()` in the functions that are increasing the wait group count.
This allows us to ensure that any background tasks that are running, will finish gracefully before the application terminates.

## 58-13. Populating the database
Copy the db.sql file and paste it in IDE's database console, then highlight the things you wanna execute and click the execute button and it will
create some tables.

Note: Whenever we're working with currency, we tend to store the info in the DB as a whole number, as an integer, and then we divide by 100 to get
the actual dollars and cents value. So we **store** 3000 and 3000/100 in DB means 30 which is in dollars and cents.

The password of **admin@example.com** user is hashed in DB and the password is **verysecret**.

The **users_plans** table is a join table. The **user_id** column is a foreign key to the users table and **plan_id** foreign key to plans table

## 59-14. Adding a data package and database models
Put data directory next to cmd directory at root of the project.


## 60-15. Implementing the loginlogout functions
Anytime you're login someone in or logout someone, it's always good to renew the token that's stored in the session.

Don't give too much info in case of err on login screen, just let people know they can't login.

In order to add sth to session, you need to register that type when we start the session. You can't just put sth iin the session willy nilly, instead
you have to register it and then put it. 

For this, in initSession() , use `gob.Register()`. If we don't doo this, the app will still compile and it'll run, but the moment someone tries to
login and does it successfully, you'll get an error saying: Can't store this in the session, don't know about the type. This is one of features of a 
strongly typed language.

Anytime we have a successfully form post, you want to redirect somewhere else, so they don't accidentally submit the form twice.

After changing the code, to see the fresh results:
```shell
make restart
```

Sending email can be expensive not just in terms of money, but in terms of processing time and delay so we do it concurrently in the background. We don't
want to site there waiting, after we created the user, while we connect to a slow mail server, we just want to fire it off in the background.