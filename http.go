package hyper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// HTTP headers as registered with IANA.
// See: https://tools.ietf.org/html/rfc7231
const (
	HeaderContentType     = "Content-Type" // RFC 7231, 3.1.1.5
	HeaderAccept          = "Accept"
	HeaderAcceptLanguage  = "Accept-Language"
	HeaderAcceptProfile   = "Accept-Profile"
	HeaderIfNoneMatch     = "If-None-Match"
	HeaderIfModifiedSince = "If-Modified-Since"
	HeaderAuthorization   = "Authorization"
)

// HTTP content types
const (
	ContentTypeHyperItem         = "application/vnd.hyper-item+json"               // https://github.com/mdemuth/hyper-item
	ContentTypeHyperItemUTF8     = "application/vnd.hyper-item+json;charset=UTF-8" // https://github.com/mdemuth/hyper-item
	ContentTypeJSON              = "application/json"                              // https://tools.ietf.org/html/rfc8259
	ContentTypeURLEncoded        = "application/x-www-form-urlencoded"             // http://www.w3.org/TR/html
	ContentTypeMultipartFormData = "multipart/form-data"                           // https://tools.ietf.org/html/rfc2388
)

// Write writes a hyper-item to the response writer with the given status code.
func Write(w http.ResponseWriter, status int, i Item) {
	w.Header().Set(HeaderContentType, ContentTypeHyperItemUTF8)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(i)
}

// WriteError writes a hyper-item representation of the error to the response writer with the given status code.
func WriteError(w http.ResponseWriter, status int, err error) {
	type errorCoder interface {
		Code() string
	}

	e := Error{}
	e.Message = err.Error()
	if errC, ok := err.(errorCoder); ok {
		e.Code = errC.Code()
	}

	Write(w, status, Item{Errors: Errors{e}})
}

const NameAction = "@action"

func ActionParameter(value string) Parameter {
	return Parameter{
		Type:  TypeHidden,
		Name:  NameAction,
		Value: value,
	}
}

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func ExtractCommand(r *http.Request) Command {
	c := MakeCommand()
	ct := r.Header.Get(HeaderContentType)
	switch ct {
	case ContentTypeURLEncoded:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return c
		}
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return c
		}
		c.Action = values.Get(NameAction)
		for n, vs := range values {
			if n == NameAction {
				continue
			}
			if len(vs) == 1 {
				c.Arguments[n] = vs[0]
			} else {
				c.Arguments[n] = vs
			}
		}
		return c
	case ContentTypeMultipartFormData:
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			return c
		}
		c.Action = r.Form.Get(NameAction)
		for n, vs := range r.Form {
			if n == NameAction {
				continue
			}
			if len(vs) == 1 {
				c.Arguments[n] = vs[0]
			} else {
				c.Arguments[n] = vs
			}
		}
		return c
	default:
		err := json.NewDecoder(r.Body).Decode(&c.Arguments)
		if err != nil {
			return c
		}
		if p := c.Arguments.String(NameAction); p != "" {
			delete(c.Arguments, NameAction)
			c.Action = p
		}
	}
	return c
}

func MakeCommand() Command {
	return Command{
		Arguments: Arguments{},
	}
}

type Command struct {
	Action    string
	Arguments Arguments
}

type Arguments map[string]interface{}

func (a Arguments) String(key string) string {
	v, ok := a[key]
	if !ok {
		return ""
	}
	switch v := v.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (a Arguments) Float64(key string) float64 {
	v, ok := a[key]
	if !ok {
		return 0
	}
	switch v := v.(type) {
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}

func (a Arguments) Int64(key string) int64 {
	v, ok := a[key]
	if !ok {
		return 0
	}
	switch v := v.(type) {
	case float64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		return 0
	}
}

func (a Arguments) Int(key string) int {
	return int(a.Int64(key))
}

func (a Arguments) Bool(key string) bool {
	v, ok := a[key]
	if !ok {
		return false
	}
	switch v := v.(type) {
	case bool:
		return v
	case float64:
		return v > 0
	case string:
		b, _ := strconv.ParseBool(v)
		return b
	default:
		return false
	}
}

func (a Arguments) Bytes(key string) []byte {
	v, ok := a[key]
	if !ok {
		return nil
	}
	switch v := v.(type) {
	case string:
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
		bs, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil
		}
		return bs
	default:
		return nil
	}
}

func ExtractContentType(h http.Header) (ContentType, error) {
	v := h.Get(HeaderContentType)
	ct := ContentType{}
	return ct, ct.Parse(v)
}

type ContentType struct {
	Type       string
	Subtype    string
	Parameters map[string]string
}

func (ct *ContentType) Parse(v string) error {
	sIndex := strings.Index(v, "/")
	if sIndex < 0 {
		return nil
	}
	ct.Type = strings.TrimSpace(v[:sIndex])
	v = v[sIndex+1:]
	pIndex := strings.Index(v, ";")
	if pIndex < 0 {
		ct.Subtype = v
		return nil
	}
	ct.Subtype = strings.TrimSpace(v[:pIndex])
	v = v[pIndex+1:]
	ps := strings.Split(v, ";")
	if len(ps) == 0 {
		return nil
	}
	if ct.Parameters == nil {
		ct.Parameters = map[string]string{}
	}
	for _, p := range ps {
		kv := strings.Split(p, "=")
		ct.Parameters[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return nil
}

func (ct ContentType) String() string {
	var buf bytes.Buffer
	buf.WriteString(ct.Type)
	buf.WriteString("/")
	buf.WriteString(ct.Subtype)
	for k, v := range ct.Parameters {
		buf.WriteString(";")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}
	return buf.String()
}

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				clearHeader(w.Header())
				WriteError(w, http.StatusInternalServerError, fmt.Errorf("recovered from: %v", err))
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func clearHeader(h http.Header) {
	for k := range h {
		h.Del(k)
	}
}
