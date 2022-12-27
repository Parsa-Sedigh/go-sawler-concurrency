package main

import (
	"context"
	"encoding/gob"
	"final-project/data"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

// create a test config for our test. This is what we use as the app for our tests.
var testApp Config

func TestMain(m *testing.M) {
	// user.Data is a non-primitive type we wanna put into the session.
	gob.Register(data.User{})

	pathToManual = "./../../pdf"
	tmpPath = "./../../tmp"

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	testApp = Config{
		Session:       session,
		Db:            nil,
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Wait:          &sync.WaitGroup{},
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
		Models:        data.TestNew(nil), // by passing nil, we don't use a DB in our tests.
	}

	// create dummy mailer for sending mails:
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	testApp.Mailer = Mail{
		Wait:       testApp.Wait,
		ErrorChan:  errorChan,
		MailerChan: mailerChan,
		DoneChan:   mailerDoneChan,
	}

	go func() {
		for {
			select {
			case <-testApp.Mailer.MailerChan: // when this runs, it means we tried to send email
				/*  Since we're never firing off the code that sends email in our tests(why? Because we''re not testing our email functionality!), so
				here we need to decrement our wait group byy 1(since the decrementing code for sending email is in it's function, but in tests, we're not
				testing that function, we need to decrement the wait group here): */
				testApp.Wait.Done()
			case <-testApp.Mailer.ErrorChan:
			case <-testApp.Mailer.DoneChan:
				return

			}
		}
	}()

	// listen for errors:
	go func() {
		for {
			select {
			case err := <-testApp.ErrorChan:
				testApp.ErrorLog.Println(err)
			case <-testApp.ErrorChanDone:
				return
			}
		}
	}()

	// this will run all of our tests after it sets up the environment
	os.Exit(m.Run())
}

/*
	A function that will add session information to our request(because we don't want to test redis but we need the session in some of the

functions we wanna test, so we just stub that session with this function)

This functions gives us a means of getting session information into and out of any request we pass it.
*/
func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)
	}

	return ctx
}
