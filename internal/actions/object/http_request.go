package object

import (
	"io"
	"net/http"
	"net/url"

	"github.com/zitadel/zitadel/internal/actions"
)

// HTTPRequestField accepts the http.Request by value, so it's not mutated
func HTTPRequestField(httpRequest http.Request) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		header := make(map[string][]string, 0)
		for k, v := range httpRequest.Header {
			header[k] = v
		}
		body := ""
		b, err := io.ReadAll(httpRequest.Body)
		if err == nil {
			body = string(b)
		}

		form := make(map[string][]string)
		postForm := make(map[string][]string)
		err = httpRequest.ParseForm()
		if err == nil {
			for k, v := range httpRequest.Form {
				form[k] = v
			}
			for k, v := range httpRequest.PostForm {
				postForm[k] = v
			}
		}

		return c.Runtime.ToValue(&httpRequestCopy{
			Method:           httpRequest.Method,
			Url:              httpRequest.URL,
			Proto:            httpRequest.Proto,
			ProtoMajor:       httpRequest.ProtoMajor,
			ProtoMinor:       httpRequest.ProtoMinor,
			Body:             body,
			ContentLength:    httpRequest.ContentLength,
			TransferEncoding: httpRequest.TransferEncoding,
			Host:             httpRequest.Host,
			Form:             form,
			PostForm:         postForm,
			Trailer:          httpRequest.Trailer,
			RemoteAddr:       httpRequest.RemoteAddr,
			RequestURI:       httpRequest.RequestURI,
		})
	}
}

type httpRequestCopy struct {
	Method           string
	Url              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Body             string
	ContentLength    int64
	TransferEncoding []string
	Host             string
	Form             map[string][]string
	PostForm         map[string][]string
	Trailer          map[string][]string
	RemoteAddr       string
	RequestURI       string
}
