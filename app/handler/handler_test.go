package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"

	"github.com/stretchr/testify/assert"
)

func TestAccountRegistration(t *testing.T) {
	c := setup(t)
	defer c.Close()

	func() {
		resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()

	func() {
		resp, err := c.Get("/v1/accounts/john")
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()
}

func setup(t *testing.T) *C {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Dao.InitAll(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

type C struct {
	App    *app.App
	Server *httptest.Server
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
