package main

import (
	"database/sql"
	"encoding/gob"
	"final-project/data"
	"fmt"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initDB()

	// create sessions
	session := initSession()

	// create loggers
	// in production, you're gonna write to a file,
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// we want to find out where the error took place, so add log.Lshortfile
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create wait group
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session:  session,
		Db:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg, // we have a wait group available to our entire application
		Models:   data.New(db),
	}

	// set up mail
	app.Mailer = app.createMail()
	go app.listenForMail()

	// listen for signals(SIGTERM and SIGINT)
	go app.listenForShutdown()

	// listen for web connections. This requires that we have sth like a routes file and also handlers
	app.serve()
}

// this function starts a web server
func (app *Config) serve() {
	// start http server
	// srv means serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // the address we're gonna listen on(webPort with any IP address on this particularR machine)
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initDB() *sql.DB {
	/* We declared this function to call another function just so that we can try to connect to the DB repeatedly if necessary(not necessary to odo this,
	it's just cleaner). */
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}

	return conn
}

/* We want to connect to DB some fixed number of times and if we can't do it after that many tries, then will just die.*/
func connectToDB() *sql.DB {
	counts := 0

	// dsn is connection string
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not ready...")
		} else {
			log.Println("connected to database")
			return connection
		}

		// if we have that error(which means we didn't return from this func), we don't want to stop at this point, we wanna try a few more times:
		if counts > 10 {
			return nil
		}

		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second) // 1 second should be enough time to to the DB
		counts++

		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// just to be safe, we ping the DB and again if there was an error, we return from the function
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {
	gob.Register(data.User{})

	// set up session
	session := scs.New()

	// with this line, we tell session store all of our info for every session in redis
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	// this actually won't be secure in localhost connection but it will be secure when it goes live
	session.Cookie.Secure = true

	return session
}

// we connect to redis using this function
func initRedis() *redis.Pool {
	// this variable is a pool of redis connections
	redisPool := &redis.Pool{
		// maximum time for an idle connection:
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return redisPool
}

func (app *Config) listenForShutdown() {
	// the size of 1 is unnecessary but we put there anyway
	quit := make(chan os.Signal, 1)

	/* When we get the interrupt signal(syscall.SIGINT) to stop the application OR syscall.SIGTERM to terminate the app which are the two things we're
	listening for, block on the quit channel. */
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	/* This just pauses until we actually get the request to interrupt or terminate the running the app*/
	<-quit

	/* If we get here, it means we had received sth from the quit channel.  */
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// perform any cleanup taSks
	app.InfoLog.Println("would run cleanup tasks...")

	/* After we run any cleanup tasks, we're gonna block until the wait group is empty and once all the semaphore count in that wait group is 0.
	we move to the next line.*/
	app.Wait.Wait()

	// After we finished sending any email(`app.Wait` wait group is empty)to tell the goroutine to quit, send a true value to done channel
	app.Mailer.DoneChan <- true

	app.InfoLog.Println("closing channels and shutting down application ...")

	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
}

func (app *Config) createMail() Mail {
	// create channels(the 3 channels we're gonna use in email)
	errorChan := make(chan error)

	/* Do we want the mailerChan as an unbuffered channel?
	If we use an unbuffered chan, it means we can only queue up one message at a time. So let's make this a buffered channel. We allow
	100 messages to be in that particular channel at any given time before it blocks.*/
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	// things like domain oor roost will come from the environment or command line flags or ... . But we just hard code them here!
	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025, // mailhog's port tis 1025
		Encryption:  "None",
		FromAddress: "info@mycompany.com",
		FromName:    "info",
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
		Wait:        app.Wait,
	}

	return m
}
