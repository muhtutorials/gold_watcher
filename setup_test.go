package main

import (
	"bytes"
	"fyne.io/fyne/v2/test"
	"gold_watcher/repository"
	"io"
	"net/http"
	"os"
	"testing"
)

var testApp Config

func TestMain(m *testing.M) {
	a := test.NewApp()
	testApp.App = a
	testApp.MainWindow = a.NewWindow("")
	testApp.HTTPClient = client
	testApp.DB = repository.NewTestRepository()
	os.Exit(m.Run())
}

var jsonToReturn = `
{
  "ts": 1691338624197,
  "tsj": 1691338615094,
  "date": "Aug 6th 2023, 12:16:55 pm NY",
  "items": [
    {
      "curr": "USD",
      "xauPrice": 1942.795,
      "xagPrice": 23.6335,
      "chgXau": 7.005,
      "chgXag": 0.032,
      "pcXau": 0.3619,
      "pcXag": 0.1356,
      "xauClose": 1935.79,
      "xagClose": 23.6015
    }
  ]
}
`

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

var client = NewTestClient(func(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
		Header:     make(http.Header),
	}
})
