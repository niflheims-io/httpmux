package httpmux

import (
	"net/http"
	"io/ioutil"
	"net/textproto"
	"html/template"
)

func newRequest(w http.ResponseWriter, r *http.Request, tpl *template.Template, ctx *context) *Request {
	resp := Response{
		w : w,
		tpl : tpl,
		gzipOn : false,
	}
	req := Request{
		r : r,
		w : &resp,
		ctx : ctx,
	}
	return &req
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

type Request struct  {
	r *http.Request
	w *Response
	params Params
	ctx *context
}

func (r *Request) HttpRequest() *http.Request {
	return r.r
}

func (r *Request) Response() *Response {
	return r.w
}

func (r *Request) Ctx() *context {
	return r.ctx
}

func (r *Request) Query(key string) string {
	if key == "" {
		return ""
	}
	return r.params.ByName(key)
}

func (r *Request) RawQuery() string {
	return r.r.URL.RawPath
}


func (r *Request) HeadGet(key string) string {
	if key == "" {
		return ""
	}
	return r.r.Header.Get(key)
}


func (r *Request) HeadSet(key string, val string) {
	if key == "" {
		return
	}
	r.r.Header.Set(key, val);
}

func (r *Request) HeadDel(key string, val string) {
	if key == "" {
		return
	}
	r.r.Header.Del(key);
}

func (r *Request) FormValue(key string) string {
	val := r.r.PostFormValue(key)
	if val == "" {
		val = r.r.PostForm.Get(key)
	}
	if val == "" {
		val = r.r.Form.Get(key)
	}
	return val
}

type UploadFile struct  {
	Data []byte
	Name string
	MIMEHeader textproto.MIMEHeader
}

func (r *Request) FormFile(key string) (*UploadFile, error) {
	parseFormErr := r.r.ParseForm()
	if parseFormErr != nil {
		return nil, parseFormErr
	}
	mFile, mFileHeader, mFileErr := r.r.FormFile(key)
	if mFileErr != nil {
		return nil, mFileErr
	}
	defer mFile.Close()
	data, readErr := ioutil.ReadAll(mFile)
	if readErr != nil {
		return nil, readErr
	}
	fileName := mFileHeader.Filename
	uFile := UploadFile{
		Data: data,
		Name: fileName,
		MIMEHeader: mFileHeader.Header,
	}
	return &uFile, nil
}
