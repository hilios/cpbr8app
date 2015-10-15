package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
)

var all = strings.Join([]string{OPTIONS, GET, POST, PUT, DELETE}, ", ")

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
		l := NewHttpLogger(*r)
		defer l.Print()
		// Add some usefull headers
		h := rw.Header()
		h.Add("Content-Type", "application/json; charset=utf-8")
		h.Add("Connection", "close")
		// CORS Headers
		h.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, "+
			"Content-Type, Accept")
		h.Add("Access-Control-Allow-Origin", "*")
		h.Add("Access-Control-Allow-Methods", "*")

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
			h.Add("Allow", all)
		}
		// Abort with a 405 status
		if handler == nil {
			Abort(&rw, http.StatusMethodNotAllowed)
			return
		}
		// Parse request data
		if err := r.ParseForm(); err != nil {
			Abort(&rw, http.StatusBadRequest)
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
		content, err := json.Marshal(data)
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
