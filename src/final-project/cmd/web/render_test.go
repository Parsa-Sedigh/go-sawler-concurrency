package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/* The AddDefaultData gets *templateData and *http.Request . So we need to build a http.Request ourselves here. */
func TestConfig_AddDefaultData(t *testing.T) {
	// the second arg doesn't matter because we're not using it.
	req, _ := http.NewRequest("GET", "/", nil)

	// add session data to the built req:
	ctx := getCtx(req)

	// with this line, this req can now accept session information(we can put things into it, takes things out of it and if sth exists)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := testApp.AddDefaultData(&TemplateData{}, req)
	if td.Flash != "flash" {
		t.Error("failed to get flash data")
	}
	if td.Warning != "warning" {
		t.Error("failed to get warning data")
	}
	if td.Error != "error" {
		t.Error("failed to get error data")
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	auth := testApp.IsAuthenticated(req)

	if auth {
		t.Error("returns true for authenticated, when it should be false")
	}

	testApp.Session.Put(ctx, "userID", 1)

	auth = testApp.IsAuthenticated(req)

	if !auth {
		t.Error("returns false for authenticated, when it should be true")
	}

}

/* The `render` function takes advantage of pathToTemplates variable(which is assumed that we're in root level of the project).
But when testing, we need to change the value fo that variable. Because when we run a test, we're not running it from the root level
of the project, we're running it wherever that test happens to live.*/

func TestConfig_render(t *testing.T) {
	pathToTemplates = "./templates"

	// response:
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.render(rr, req, "home.page.gohtml", &TemplateData{})

	if rr.Code != 200 {
		t.Error("failed too render page")
	}
}
