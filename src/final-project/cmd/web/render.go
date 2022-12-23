package main

import (
	"final-project/data"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	// a map of strings
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64

	// we can put anything we want in here:
	Data map[string]any

	// these are alerts oor similar stuff
	Flash   string
	Warning string
	Error   string

	/* is user is authenticated: */
	Authenticated bool
	Now           time.Time
	User          *data.User
}

// t is the name oof the template to render
/* since there might be the case where we don't want to send any data to a template we wanna render, we made the third arg a pointer so that func init() {
could be nil instead oof a  zero--value and we can check for nil easily ini the func.*/
func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	// these have tto be here for every template that we want to render
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	// actual template that we wanna render:
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	if td == nil {
		// with this, td is now a zero-value of TemplateData, at least it's not nil, because we can't doo anything with nil
		td = &TemplateData{}
	}

	/// parse the templates
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // we can't go any further, so return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	/* PopString makes sure that as soon as this data is read, it's removed from the session, which is convenient for messages that you only
	want to display once.*/
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		td.Authenticated = true

		// cast the found user to data.user type:
		user, ok := app.Session.Get(r.Context(), "user").(data.User)
		if !ok {
			app.ErrorLog.Println("can't get user from the session")
		} else {
			// add all the info of the user in the template(it needs to be a pointer)
			td.User = &user
		}
	}

	td.Now = time.Now()

	return td
}

func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
