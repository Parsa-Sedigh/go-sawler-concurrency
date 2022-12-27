package main

import (
	"final-project/data"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

/* Instead of just doing one test per page, we'll write a simple table test. In other words, we'll set up some data called pageTests */
var pageTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	handler            http.HandlerFunc
	sessionData        map[string]any
	expectedHTML       string // once we render the page, we can look at the html generated and search for a particular string
}{{
	name:               "home",
	url:                "/",
	expectedStatusCode: http.StatusOK,
	handler:            testApp.HomePage,
},
	{
		name:               "login page",
		url:                "/login",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:               "logout page",
		url:                "/logout",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.LoginPage,

		// we should be able to logout only when we're logged in, soo we should see some session data
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
}

func Test_Pages(t *testing.T) {
	// overwrite pathToTemplates var and it should be now relative to current file:
	pathToTemplates = "./templates"

	for _, e := range pageTests {
		// substitutes or takes the place of the response writer:
		rr := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(e.sessionData) > 0 {
			for key, value := range e.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		e.handler.ServeHTTP(rr, req)

		// perform the test:
		if rr.Code != http.StatusOK {
			t.Errorf("%s failed: expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if len(e.expectedHTML) > 0 {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("%s failed: expected to ofind %s, but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func TestConfig_PostLoginPage(t *testing.T) {
	/* First we need to overwrite the value of pathToTemplates var so that the render function can find the appropriate page if necessary. */
	pathToTemplates = "./templates"

	postedData := url.Values{
		"email":    {"admin@example.com"},
		"password": {"abc123abc123abc123abc123abc123"},
	}

	// rr stands for response recorder
	rr := httptest.NewRecorder()

	// second param is not used, but it's nice to have it match what it's actually going to
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(testApp.PostLoginPage)
	handler.ServeHTTP(rr, req)

	// in PostLoginPage, we're supposed to get back http.StatusSeeOther. We can test this!
	if rr.Code != http.StatusSeeOther {
		t.Error("wrong code returned")
	}

	if !testApp.Session.Exists(ctx, "userID") {
		t.Error("did not find userID in the session")
	}
}

func TestConfig_SubscribeToPlan(t *testing.T) {
	rr := httptest.NewRecorder()

	// in this test, we're actually using the second param off NewRequest()
	req, _ := http.NewRequest("GET", "/subscribe?id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// this particular handler that we're trying to test here, requires that there be a user in the session, so we have to put one:
	testApp.Session.Put(ctx, "user", data.User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Active:    1,
	})

	handler := http.HandlerFunc(testApp.SubscribeToPlan)
	handler.ServeHTTP(rr, req)

	testApp.Wait.Wait()

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status code of statusSeeOther, but got %d", rr.Code)
	}
}
