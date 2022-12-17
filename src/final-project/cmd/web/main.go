package main

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"time"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initDB()

	// create sessions

	// create channels

	// create waitgroup

	// set up the application config

	// set up mail

	// listen for web connections
}

func initDB() *sql.DB {
	/* We declared this function to call another function just so that we can try to connect to the DB repeatedly if necessary(not necessary to odo this,
	it's just cleaner). */
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}
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
		time.Sleep(1 * time.Second) // 1 second should be enough time to too the DB
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
