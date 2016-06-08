package httpmux

import (
	"net/http"
	"io/ioutil"
	"net/textproto"
)

type request struct  {
	httpRequest *http.Request
	httpResponse http.Response
	params map[string]string
}

func (self *request) setParam(key, val string) {
	if key == "" {
		return ""
	}
	self.params[key] = val
}

func (self *request) Query(key string) string {
	if key == "" {
		return ""
	}
	if val, get := self.params[key]; get {
		return val
	}
	return ""
}

func (self *request) RawQuery() string {
	return self.httpRequest.URL.RawPath
}

func (self *request) HttpRequest() *http.Request {
	return self.httpRequest
}

func (self *request) HeadGet(key string) string {
	if key == "" {
		return ""
	}
	return self.httpRequest.Header.Get(key)
}


func (self *request) HeadSet(key string, val string) {
	if key == "" {
		return
	}
	self.httpRequest.Header.Set(key, val);
}

func (self *request) HeadDel(key string, val string) {
	if key == "" {
		return
	}
	self.httpRequest.Header.Del(key, val);
}

func (self *request) FormValue(key string) string {
	val := self.httpRequest.PostFormValue(key)
	if val == "" {
		val = self.httpRequest.PostForm.Get(key)
	}
	if val == "" {
		val = self.httpRequest.Form.Get(key)
	}
	return val
}

type UploadFile struct  {
	Data []byte
	Name string
	MIMEHeader textproto.MIMEHeader
}

func (self *request) FormFile(key string) (*UploadFile, error) {
	parseFormErr := self.httpRequest.ParseForm()
	if parseFormErr != nil {
		return nil, parseFormErr
	}
	mFile, mFileHeader, mFileErr := self.httpRequest.FormFile(key)
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

func newContext(w http.ResponseWriter, r *http.Request, res *resMap) *Context {
	return &Context{Request:r, Response:w, res:res}
}

type Context struct  {
	Request *request
	Response *http.Response
	res *resMap
}
