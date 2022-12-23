package main

import "net/http"

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
	// create a user

	// send an activation email

	// subscribe the user to an account
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url(we want to have a secure link in the email)

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached
}
