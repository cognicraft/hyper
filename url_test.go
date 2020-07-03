package hyper

import (
	"crypto/tls"
	"net/http"
	"strings"
	"testing"
)

func TestExternalURLResolver(t *testing.T) {
	tests := []struct {
		req       *http.Request
		relFmt    string
		relParams []interface{}
		out       string
	}{
		{
			req:    mustRequest("GET", "/"),
			relFmt: ".",
			out:    "/",
		},
		{
			req:    mustRequest("GET", "http://localhost:8080/foo/bar/"),
			relFmt: ".",
			out:    "http://localhost:8080/foo/bar/",
		},
		{
			req:    mustRequest("GET", "http://localhost:8080/foo/bar/"),
			relFmt: "..",
			out:    "http://localhost:8080/foo/",
		},
		{
			req:    mustRequest("GET", "http://localhost:8080/foo/bar/"),
			relFmt: "../test/search?q=fizzle",
			out:    "http://localhost:8080/foo/test/search?q=fizzle",
		},
		{
			req:       mustRequest("GET", "http://localhost:8080/foo/bar/"),
			relFmt:    "../%s/search?q=fizzle",
			relParams: []interface{}{"dizzle"},
			out:       "http://localhost:8080/foo/dizzle/search?q=fizzle",
		},
		{
			req:    mustRequest("GET", "http://localhost:8080/foo/bar/search?q=fizzle"),
			relFmt: "",
			out:    "http://localhost:8080/foo/bar/search?q=fizzle",
		},
		{
			req:    mustRequest("GET", "https://localhost:8080/foo/bar/"),
			relFmt: "../test/search?q=fizzle",
			out:    "https://localhost:8080/foo/test/search?q=fizzle",
		},
		{
			req: mustRequest("GET", "http://localhost:8080/", h{HeaderXForwardedProto, "http"}),
			out: "http://localhost:8080/",
		},
		{
			req: mustRequest("GET", "http://localhost:8080/", h{HeaderXForwardedProto, "https"}),
			out: "https://localhost:8080/",
		},
		{
			req: mustRequest("GET", "http://localhost:8080/foo/bar/flash",
				h{HeaderXForwardedProto, "https"},
				h{HeaderXForwardedHost, "10.0.0.1:1234"},
			),
			out: "https://10.0.0.1:1234/foo/bar/flash",
		},
		{
			req: mustRequest("GET", "http://localhost:8080/foo/bar/flash",
				h{HeaderXForwardedProto, "https"},
				h{HeaderXForwardedHost, "10.0.0.1:1234"},
			),
			out: "https://10.0.0.1:1234/foo/bar/flash",
		},
		{
			req: mustRequest("GET", "http://localhost:8080/foo/bar/flash",
				h{HeaderXForwardedProto, "https"},
				h{HeaderXForwardedHost, "10.0.0.1:1234"},
				h{HeaderXForwardedPath, "/wizzle/dizzle?q=sub"},
			),
			out: "https://10.0.0.1:1234/wizzle/dizzle?q=sub",
		},
	}
	for _, test := range tests {
		resolve := ExternalURLResolver(test.req)
		got := resolve(test.relFmt, test.relParams...).String()
		if test.out != got {
			t.Errorf("\nwant:\n%s\n got:\n%s", test.out, got)
		}
	}
}

func mustRequest(method string, url string, headers ...h) *http.Request {
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(url, "https") {
		r.TLS = &tls.ConnectionState{}
	}
	for _, h := range headers {
		r.Header.Set(h.Key, h.Value)
	}
	return r
}

type h struct {
	Key   string
	Value string
}
