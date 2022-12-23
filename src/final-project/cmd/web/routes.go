package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *Config) routes() http.Handler {
	// create router

	// mux means multiplexer
	mux := chi.NewRouter()

	// set up middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	// define application routes
	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate-account", app.ActivateAccount)
	mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain: "localhost",

			// we'll be sending test mails to mailshot which will be running in our docker images
			Host:        "localhost",
			Port:        1025,   // mailhog's port tis 1025
			Encryption:  "None", // in development
			FromAddress: "info@mycompany.com",
			FromName:    "info",
			ErrorChan:   make(chan error),
		}

		msg := Message{
			To:      "me@here.com",
			Subject: "Test email",
			Data:    "Hello, world.",
		}

		m.sendMail(msg, make(chan error))

	})

	return mux
}
