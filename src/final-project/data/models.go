package data

import (
	"database/sql"
	"time"
)

// if you can't do your DB operation in 3 seconds, sth has gone wrong
const dbTimeout = time.Second * 3

// our connection pool. This is what we use to connect to DB
var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: &User{},
		Plan: &Plan{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User UserInterface
	Plan PlanInterface
}
