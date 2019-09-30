package hyper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	HeaderXForwardedHost   = "X-Forwarded-Host"
	HeaderXForwardedFor    = "X-Forwarded-For"
	HeaderXForwardedProto  = "X-Forwarded-Proto"
	HeaderXForwardedScheme = "X-Forwarded-Scheme"
	HeaderXForwardedPath   = "X-Forwarded-Path"
	HeaderXRealIP          = "X-Real-IP"
	HeaderXCorrelationID   = "X-Correlation-ID"
	HeaderForwarded        = "Forwarded"
)

func ExternalScheme(r *http.Request) string {
	scheme := "http"
	if isHTTPS := r.TLS != nil; isHTTPS {
		scheme = "https"
	}
	if forwardedProto := r.Header.Get(HeaderXForwardedProto); forwardedProto != "" {
		scheme = forwardedProto
	}
	return scheme
}

func ExternalHost(r *http.Request) string {
	host := r.Host
	if forwardedHost := r.Header.Get(HeaderXForwardedHost); forwardedHost != "" {
		host = forwardedHost
	}
	return host
}

func ExternalPath(r *http.Request) string {
	path := r.URL.Path
	if forwardedPath := r.Header.Get(HeaderXForwardedPath); forwardedPath != "" {
		path = forwardedPath
	}
	return path
}

func ExternalURL(r *http.Request) *url.URL {
	scheme := ExternalScheme(r)
	host := ExternalHost(r)
	path := ExternalPath(r)
	query := r.URL.RawQuery

	exURL := fmt.Sprintf("%s://%s%s", scheme, host, path)
	if len(query) > 0 {
		exURL = fmt.Sprintf("%s://%s%s?%s", scheme, host, path, query)
	}

	result, _ := url.Parse(exURL)
	return result
}

type Resolver interface {
	Resolve(format string, args ...interface{}) *url.URL
}

type ResolverFunc func(format string, args ...interface{}) *url.URL

func (fn ResolverFunc) ResolverFunc(format string, args ...interface{}) ResolverFunc {
	return NewURLResolver(fn.Resolve(format, args...)).Resolve
}

func (fn ResolverFunc) Resolve(format string, args ...interface{}) *url.URL {
	return fn(format, args...)
}

func ExternalURLResolver(r *http.Request) ResolverFunc {
	return NewURLResolver(ExternalURL(r))
}

func ResolveURL(baseURL *url.URL, format string, args ...interface{}) *url.URL {
	if baseURL == nil {
		baseURL, _ = url.Parse("/")
	}
	rel := fmt.Sprintf(format, args...)
	res, _ := baseURL.Parse(rel)
	return res
}

func NewURLResolver(baseURL *url.URL) ResolverFunc {
	return func(format string, args ...interface{}) *url.URL {
		return ResolveURL(baseURL, format, args...)
	}
}

func ExtractRemote(r *http.Request) string {
	if forwardedFor := r.Header.Get(HeaderXForwardedFor); forwardedFor != "" {
		return forwardedFor
	}
	remParts := strings.Split(r.RemoteAddr, ":")
	if len(remParts) > 0 {
		return remParts[0]
	}
	return ""
}
