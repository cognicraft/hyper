package hyper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	HeaderLocation        = "Location"
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
	switch {
	case strings.HasPrefix(ct, ContentTypeURLEncoded):
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
	case strings.HasPrefix(ct, ContentTypeMultipartFormData):
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			return c
		}
		for n, vs := range r.MultipartForm.Value {
			if n == NameAction && len(vs) > 0 {
				c.Action = vs[0]
			}
			if len(vs) == 1 {
				c.Arguments[n] = vs[0]
			} else {
				c.Arguments[n] = vs
			}
		}
		for n, vs := range r.MultipartForm.File {
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
		if strings.HasPrefix(v, "data:") {
			if i := strings.Index(v, ","); i > 0 {
				// drop prefix
				uv, err := unescape(v[i+1:])
				if err != nil {
					return nil
				}
				v = string(uv)
			}
		}
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
		bs, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Printf("%v", err)
			return nil
		}
		return bs
	case *multipart.FileHeader:
		f, err := v.Open()
		if err != nil {
			return nil
		}
		bs, err := ioutil.ReadAll(f)
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

func unescape(s string) ([]byte, error) {
	var buf = new(bytes.Buffer)
	reader := strings.NewReader(s)

	for {
		r, size, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if size > 1 {
			return nil, fmt.Errorf("rfc2396: non-ASCII char detected")
		}

		switch r {
		case '%':
			eb1, err := reader.ReadByte()
			if err == io.EOF {
				return nil, fmt.Errorf("rfc2396: unexpected end of unescape sequence")
			}
			if err != nil {
				return nil, err
			}
			if !isHex(eb1) {
				return nil, fmt.Errorf("rfc2396: invalid char 0x%x in unescape sequence", r)
			}
			eb0, err := reader.ReadByte()
			if err == io.EOF {
				return nil, fmt.Errorf("rfc2396: unexpected end of unescape sequence")
			}
			if err != nil {
				return nil, err
			}
			if !isHex(eb0) {
				return nil, fmt.Errorf("rfc2396: invalid char 0x%x in unescape sequence", r)
			}
			buf.WriteByte(unhex(eb0) + unhex(eb1)*16)
		default:
			buf.WriteByte(byte(r))
		}
	}
	return buf.Bytes(), nil
}

func isHex(c byte) bool {
	switch {
	case c >= 'a' && c <= 'f':
		return true
	case c >= 'A' && c <= 'F':
		return true
	case c >= '0' && c <= '9':
		return true
	}
	return false
}

// borrowed from net/url/url.go
func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
