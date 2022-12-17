# Section 6: 6. Final Project - Building a Subscription Service

## 46-1. What we'll cover in this section

## 47-2. Setting up a simple web application
In the main func, since we're sending mails, we need to use wait groups, because we're sending emails concurrently. Now if someone decides
to restart the application for whatever reason(maybe doing some maintenance or ...), you don't want to just say: Stop the application. Instead,
we want to wait until any mail that is waiting to be sent, gets delivered and then you stop the app and this is an ideal situation for wait group.

There are a couple of different drives for postgres, one is PQ but it says use the jackc's one instead which is named `pgconn`! So let's use that one and then
`pgx/v4`. For session management we use `scs/v2` . This package allows different stores for session data. You can store it in the database, in a cookie, but we 
want too use redis and for this, install: `scs/redisstore`. Also install the chi router: `go-chi/chi/v5`.

## 48-3. Setting up our Docker development environment
We're gonna use redis as a store for our session info and we're also gonna require some kind of mail server, a dummy mail server that will capture mail
for us and we're gonna use docker for that.

**Mailhog**(dummy mail server): it allows us to send email without actually sending email! Mailhog captures it and if offers a web interface where we can go and look 
at our email.

In docker-compose file, we specified a volume for storing the data for postgres in the `./db-data/postgresa/...`. Now in a lot of cases, 
docker will create that for you, but for safety, you can create a new folder at the root level of your project called db-data and `postgres` directory and
one for `redis`.

Then you scan run:
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
## 51-6. Adding sessions & Redis
## 52-7. Setting up the application config
## 53-8. Setting up a route & handler for the home page, and starting the web server
## 54-9. Setting up templates and building a render function
## 55-10. Adding session middleware
## 56-11. Setting up additional stub handlers and routes
## 57-12. Implementing graceful shutdown
## 58-13. Populating the database
## 59-14. Adding a data package and database models
## 60-15. Implementing the loginlogout functions
