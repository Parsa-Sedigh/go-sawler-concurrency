package main

import (
	"final-project/data"
	"fmt"
	"html/template"
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context()) // with this, the session token is renewed

	// parse form post
	err := r.ParseForm()
	if err != nil {
		// in a real app, you can redirect user to that screen with an error message
		app.ErrorLog.Println(err)
	}

	// get the email and password from form post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Till here, we found the user, but let's check the password
	/* PasswordMatches will take the password that user sent us, compare the hash that we have in the DB with the password that supplied,
	if they match, validPassword is true, if they do not, then the user entered the wrong password. */
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		// note: We don't want to give any more info than "Invalid credentials" away for security reasons
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		/* send an email notification to notify user he has an invalid login in his account.
		Note: In production you want to keep track of the number of logins and onplay send it after the third check or sth like that.*/
		msg := Message{
			// we'll use the default values for FromName and From fields empty(we have set some defaults for them in createMail func).
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Invalid login attempt!",
		}

		// this will send the email in background
		app.sendEmail(msg)

		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	/* okay, so log user in. How do we log a user in? We just put user id in the session.  */
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	/* After seeing the flash message and reloading, that should be gone because that's pulled out of the session using Session.Pop()  */
	app.Session.Put(r.Context(), "flash", "Successful login")

	// redirect the user
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())

	// once user is logged out,	it's always good practice to renew the session token
	_ = app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	// we know that we will handling a form POST and in go anytime you get a form POST, the first thing you wanna do is to parse form data
	err := r.ParseForm()
	if err != nil { // it's rare that this error takes place!
		app.ErrorLog.Println(err)
	}

	/* TODO - validate data: For example we wanna make sure that this user is not already registered. We wanna make sure that they filled all the necessary
	stuff.*/

	// create a user
	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),

		/* even though we don't need to specify this(zero-value default), for readability, let's set Active to 0(we're creating a user that is inactive for now).
		So even though these would get the default values that are correct, this just makes it clear that we're creating ga user that is not an admin
		and not active.*/
		Active:  0,
		IsAdmin: 0,
	}

	// insert the user in DB
	_, err = u.Insert(u)
	if err != nil {
		// set the key to error so it shows up as an error oon the page
		app.Session.Put(r.Context(), "error", "Unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activation email
	/* We hard coded the url in this case. Normally you would read that from an environment variable or.env file. */
	url := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL), // cast signedURL to template.HTML()
	}

	// this is gonna take place in the background
	app.sendEmail(msg)

	// put a success message:
	app.Session.Put(r.Context(), "flash", "Confirmation email sent. Check your email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url(we want to have a secure link in the email)
	url := r.RequestURI

	// build a test string(sth that can be used to validate the url)
	testURL := fmt.Sprintf("http://localhost%s", url)
	okay := VerifyToken(testURL)

	// we have an invalid url signature
	if !okay {
		app.Session.Put(r.Context(), "error", "Invalid token.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate the account by 1) get the user from DB 2) set active to 1 3) save the user
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update user.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Account activated. You can now log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached

	// subscribe the user to an account
}

func (app *Config) chooseSubscription(w http.ResponseWriter, r *http.Request) {
	if !app.Session.Exists(r.Context(), "userID") {
		app.Session.Put(r.Context(), "warning", "You must log in to see this page!")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		// you might want to display an error page in this case
		app.ErrorLog.Println(err)
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}
