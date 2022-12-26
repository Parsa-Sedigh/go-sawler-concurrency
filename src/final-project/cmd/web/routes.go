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
	mux.Get("/activate", app.ActivateAccount)
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

	// this means, /plans is now /members/plans (any route handler in the below router, has an append /members to it)
	mux.Mount("/members", app.authRouter())

	return mux
}

func (app *Config) authRouter() http.Handler {
	// create a new router:
	// any route that the user should be authenticated to use it, is gonna get into this router
	mux := chi.NewRouter()

	// for this particular router(named mux), we use the Auth middleware:
	mux.Use(app.Auth)

	mux.Get("/plans", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubscribeToPlan)

	return mux
}
