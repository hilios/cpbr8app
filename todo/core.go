package todo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type DefaultResponseData struct {
	StatusInt int    `json:status_int`
	Status    string `json:status`
}

type Response struct {
	StatusInt int
	Data      interface{}
}

func (r *Response) GetStatus() int {
	return r.StatusInt
}

func (r *Response) GetStatusText() string {
	return fmt.Sprintf("%d %s", r.StatusInt, http.StatusText(r.StatusInt))
}

func (r *Response) GetData() interface{} {
	if r.Data != nil {
		return r.Data
	}
	return DefaultResponseData{
		StatusInt: r.GetStatus(),
		Status:    r.GetStatusText(),
	}

}

type Resource interface {
	Get(values url.Values) Response
	Post(values url.Values) Response
	Put(values url.Values) Response
	Delete(values url.Values) Response
}

type (
	GetNotAllowed    struct{}
	PostNotAllowed   struct{}
	PutNotAllowed    struct{}
	DeleteNotAllowed struct{}
)

func (GetNotAllowed) Get(values url.Values) Response {
	return Response{http.StatusMethodNotAllowed, nil}
}

func (PostNotAllowed) Post(values url.Values) Response {
	return Response{http.StatusMethodNotAllowed, nil}
}

func (PutNotAllowed) Put(values url.Values) Response {
	return Response{http.StatusMethodNotAllowed, nil}
}

func (DeleteNotAllowed) Delete(values url.Values) Response {
	return Response{http.StatusMethodNotAllowed, nil}
}

type Logger struct {
	initTime time.Time
	request  *http.Request
}

func (l *Logger) Print() {
	log.Printf("%s \t %s %s",
		time.Since(l.initTime).String(),
		l.request.Method,
		l.request.URL.Path)
}

func NewLogger(r *http.Request) *Logger {
	l := new(Logger)
	l.initTime = time.Now()
	l.request = r
	return l
}

func Abort(w http.ResponseWriter, code int, err error) {
	errHtml := fmt.Sprintf("<h1>%s</h1><p>%s</p>", http.StatusText(code), err)
	w.WriteHeader(code)
	w.Write([]byte(errHtml))
}

func ResponseError() Response {
	return Response{http.StatusInternalServerError, nil}
}

func RestHandler(r Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var response Response
		// Populate request.Form
		request.ParseForm()
		values := request.Form
		method := request.Method
		// Output
		logger := NewLogger(request)
		// Dispatch the method at the Resource
		switch method {
		case GET:
			response = r.Get(values)
		case POST:
			response = r.Post(values)
		case PUT:
			response = r.Put(values)
		case DELETE:
			response = r.Delete(values)
		default:
			response = Response{http.StatusMethodNotAllowed, nil}
		}
		// Encode data as a JSON
		jsonData, err := json.Marshal(response.GetData())
		if err != nil {
			http.Error(w, "JSON encoding failed",
				http.StatusInternalServerError)
		} else {
			// Output
			w.WriteHeader(response.GetStatus())
			w.Write(jsonData)
		}
		// Print a log message
		logger.Print()
	}
}
