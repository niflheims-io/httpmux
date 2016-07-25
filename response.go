package httpmux

import (
	"net/http"
	"encoding/json"
	"compress/gzip"
	"encoding/xml"
	"html/template"
)

type Response struct  {
	w http.ResponseWriter
	tpl *template.Template
	gzipOn bool
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func (r *Response) Gzip() *Response {
	r.w.Header().Add("Content-Encoding", "gzip")
	r.gzipOn = true
	return r
}

func (r *Response) Status(status int) *Response {
	r.w.WriteHeader(status)
	return  r
}

func (r *Response) Write(data []byte, contentTypes []string) error {
	writeContentType(r.w, contentTypes)
	if r.gzipOn {
		w := gzip.NewWriter(r.w)
		defer w.Close()
		defer w.Flush()
		_, err := w.Write(data)
		return err
	}
	_, err := r.w.Write(data)
	return err
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func (r *Response) String(o string) error {
	writeContentType(r.w, plainContentType)
	if r.gzipOn {
		w := gzip.NewWriter(r.w)
		defer w.Close()
		defer w.Flush()
		_, err := w.Write([]byte(o))
		return err
	}
	_, err := r.w.Write([]byte(o))
	return err
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (r *Response) Json(o interface{}) error {
	writeContentType(r.w, jsonContentType)
	if r.gzipOn {
		w := gzip.NewWriter(r.w)
		defer w.Close()
		defer w.Flush()
		err := json.NewEncoder(w).Encode(o)
		return err
	}
	return json.NewEncoder(r.w).Encode(o)
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

func (r *Response) Xml(o interface{}) error {
	writeContentType(r.w, xmlContentType)
	if r.gzipOn {
		w := gzip.NewWriter(r.w)
		defer w.Close()
		defer w.Flush()
		err := xml.NewEncoder(w).Encode(o)
		return err
	}
	return xml.NewEncoder(r.w).Encode(o)
}

var htmlContentType = []string{"text/html; charset=utf-8"}

func (r *Response) Html(tplName string, o interface{}) error {
	writeContentType(r.w, htmlContentType)
	if r.gzipOn {
		w := gzip.NewWriter(r.w)
		defer w.Close()
		defer w.Flush()
		if tplName == "" {
			return r.tpl.Execute(w, o)
		}
		return r.tpl.ExecuteTemplate(w, tplName, o)
	}
	if tplName == "" {
		return r.tpl.Execute(r.w, o)
	}
	return r.tpl.ExecuteTemplate(r.w, tplName, o)
}

