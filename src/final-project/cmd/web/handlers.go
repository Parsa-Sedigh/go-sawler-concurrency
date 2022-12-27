package main

import (
	"final-project/data"
	"fmt"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var pathToManual = "./pdf"
var tmpPath = "./tmp"

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
	if they match, validPassword is true, if they do not, then the user entered the wrong password.

	Instead of writing: user.PasswordMatches(password), use app.Models.User.... . Why? Look at the 89-8 pamphlet. */
	validPassword, err := app.Models.User.PasswordMatches(password)
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
	_, err = app.Models.User.Insert(u)
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
	err = app.Models.User.Update(*u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update user.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Account activated. You can now log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// send an email with attachments

}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	// commented this out because it's handled by the app.auth() middleware:
	//if !app.Session.Exists(r.Context(), "userID") {
	//	app.Session.Put(r.Context(), "warning", "You must log in to see this page!")
	//	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	//	return
	//}

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

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// get the id of the plan that is chosen
	id := r.URL.Query().Get("id")
	planID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorLog.Println("Error getting planID: ", err)
	}

	// get the plan from the database
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find plan.")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	// get the user from the session
	user, ok := app.Session.Get(r.Context(), "user").(data.User)

	/* if we can't get the user, chances are they're not logged in(the scenario is they went to the /plans page while they were logged in and then after some time,
	their session timed out. This isi the only way they could have this error here.) */
	if !ok {
		app.Session.Put(r.Context(), "error", "Log in first!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// generate an invoice and email it
	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		// generate the invoice
		invoice, err := app.getInvoice(user, plan)
		if err != nil {
			/* What do we do here? This function that has this error is running in the background! We have an error but we have no means oof
			sending that error back. We need to send it to a channel! But which channel? Let's create some in Config struct.*/
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your invoice",
			Data:     invoice,
			Template: "invoice",
		}

		app.sendEmail(msg)
	}()

	// generate a manual
	app.Wait.Add(1)
	// we could put all this logic into the above goroutine, but we assume that getInvoice() has a lot of logic in it.
	go func() {
		defer app.Wait.Done()

		pdf := app.generateManual(user, plan)

		/* write pdf to a temporary folder at the root level of our app.
		We named the file like below so that we know that we're never gonna have one user generating two manuals at the same instant. So this way we know
		that  by prepending the user id at the beginning oof this file name, it won't overwrite somebody else's file.*/
		err := pdf.OutputFileAndClose(fmt.Sprintf("%s/%d_manual.pdf", tmpPath, user.ID))
		if err != nil {
			app.ErrorChan <- err
			return
		}

		// send an email with the manual attached
		msg := Message{
			To:      user.Email,
			Subject: "Your manual",
			Data:    "Your user manual is attached",

			/* we want a readable name for the file when the customer downloads it. Now we could split the filename on the underscore and hope for the best,
			but let's add a new field to Message struct named AttachmentMap. Because sometimes you want to attach sth and overwrite the name.*/
			AttachmentMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("%s/%d_manual.pdf", tmpPath, user.ID),
			}, // an attachment with a custom name
		}

		app.sendEmail(msg)

		// test app error chan:
		app.ErrorChan <- err
	}()

	// subscribe the user to a plan
	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error subscribing to plan	!")
		http.Redirect(w, r, "/members/plan", http.StatusSeeOther)
		return
	}

	/* At this point, we have subscribed this user to other plan, but remember, we have user in the session and the user in the session is still
	subscribed to the old plan. So we need to update the user in the session, we just need a fresh copy from the DB (note that user.ID didn't change so it's safe to
	use it here).*/
	u, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error getting user from database!")
		http.Redirect(w, r, "/members/plan", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "user", u)

	// redirect
	app.Session.Put(r.Context(), "flash", "Subscribed!")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}

func (app *Config) getInvoice(u data.User, plan *data.Plan) (string, error) {
	// some heavy lifting in a production app...

	return plan.PlanAmountFormatted, nil
}

func (app *Config) generateManual(u data.User, plan *data.Plan) *gofpdf.Fpdf {
	// define a PDF, set it's size and margins and ...
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	// import an existing PDF:
	importer := gofpdi.NewImporter()

	// simulate the amount of time might take to create a PDF:
	time.Sleep(5 * time.Second)

	t := importer.ImportPage(pdf, fmt.Sprintf("%s/manual.pdf", pathToManual), 1, "/MediaBox")
	pdf.AddPage()

	// use the imported template for that page:
	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	// you do this(get these numbers) by getting a ruler and measuring where you want sth to appear oon the page
	pdf.SetX(75)
	pdf.SetY(150)

	// we chose a font that we're sure that is installed on every computer out there:
	pdf.SetFont("Arial", "", 12)

	// we want to write a cell that may span multiple lines:
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName, u.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}
