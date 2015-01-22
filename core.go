package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
)

type GetMethod interface {
	Get(url.Values) (int, interface{})
}

type PostMethod interface {
	Post(url.Values) (int, interface{})
}

type PutMethod interface {
	Put(url.Values) (int, interface{})
}

type DeleteMethod interface {
	Delete(url.Values) (int, interface{})
}

// Dispatches an HTTP error with the given code
func Abort(rw *http.ResponseWriter, code int) {
	err := fmt.Sprintf("%d %s", code, http.StatusText(code))
	http.Error(*rw, err, code)
}

// Reflects the different HTTP verbs into the given interface, allows the
// following methods: `Get`, `Post`, `Put`, `Delete`.
//
// Based on a Doug Black code:
// https://github.com/dougblack/sleepy/blob/master/core.go
func RestController(c interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log := NewHttpLogger(*r)
		defer log.Print()
		// Add some usefull headers
		h := rw.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "*")
		h.Set("Allow", "*")
		h.Set("Connection", "close")
		// Parse sent data
		if r.ParseForm() != nil {
			Abort(&rw, http.StatusBadRequest)
			return
		}

		var handler func(url.Values) (int, interface{})

		switch r.Method {
		case GET:
			if c, ok := c.(GetMethod); ok {
				handler = c.Get
			}
		case POST:
			if c, ok := c.(PostMethod); ok {
				handler = c.Post
			}
		case PUT:
			if c, ok := c.(PutMethod); ok {
				handler = c.Put
			}
		case DELETE:
			if c, ok := c.(DeleteMethod); ok {
				handler = c.Delete
			}
		case OPTIONS:
			handler = func(_ url.Values) (int, interface{}) {
				return http.StatusOK, ""
			}
		}
		// Abort with a 405 status
		if handler == nil {
			Abort(&rw, http.StatusMethodNotAllowed)
			return
		}
		// Create the params from GET and POST values
		params, _ := url.ParseQuery(fmt.Sprintf("%s&%s",
			r.URL.Query().Encode(), r.Form.Encode()))
		// Call the handler
		code, data := handler(params)
		// Set a default
		if data == nil {
			data = map[string]interface{}{
				"status": code,
				"text":   http.StatusText(code),
			}
		}
		// Encode
		content, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			Abort(&rw, http.StatusInternalServerError)
			return
		}
		// Write to response
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

// Simple HTTP request logger
type HttpLogger struct {
	initTime time.Time
	request  http.Request
}

// Print with the method, url and request time
func (l *HttpLogger) Print() {
	log.Printf("%s \t %s %s",
		time.Since(l.initTime).String(),
		l.request.Method,
		l.request.URL.Path)
}

// Return a new instance of the HTTP logger
func NewHttpLogger(r http.Request) *HttpLogger {
	l := new(HttpLogger)
	l.initTime = time.Now()
	l.request = r
	return l
}
