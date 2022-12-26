package main

import (
	"database/sql"
	"final-project/data"
	"github.com/alexedwards/scs/v2"
	"log"
	"sync"
)

type Config struct {
	Session  *scs.SessionManager
	Db       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
	Mailer   Mail

	// centralized channels for error handling
	ErrorChan     chan error
	ErrorChanDone chan bool
}
