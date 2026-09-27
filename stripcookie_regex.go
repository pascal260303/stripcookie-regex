// Package stripcookie a plugin to strip cookies.
package stripcookie_regex

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Cookies       []string `json:"cookies,omitempty"`
	CookieRegexes []string `json:"cookieRegexes,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// CookieStrip a CookieStrip plugin.
type CookieStrip struct {
	next          http.Handler
	cookies       []string
	cookieRegexes []*regexp.Regexp
	name          string
	splitRegexp   *regexp.Regexp
}

// New created a new CookieStrip plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Cookies) == 0 && len(config.CookieRegexes) == 0 {
		return nil, fmt.Errorf("cookies and cookieRegexes cannot both be empty")
	}

	regexes, err := compileRegexes(config.CookieRegexes)
	if err != nil {
		return nil, err
	}

	return &CookieStrip{
		cookies:       config.Cookies,
		cookieRegexes: regexes,
		next:          next,
		name:          name,
		splitRegexp:   regexp.MustCompile(` *([^=;]+?) *=[^;]+`),
	}, nil
}

func (c *CookieStrip) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cookieHeaders := req.Header.Values("cookie")
	req.Header.Del("cookie")
	for _, cookieHeader := range cookieHeaders {
		cookies := c.splitRegexp.FindAllStringSubmatch(cookieHeader, -1)
		var keep []string
		for _, cookie := range cookies {
			if !c.shouldStrip(cookie[1]) {
				keep = append(keep, cookie[0])
			}
		}
		if len(keep) > 0 {
			req.Header.Add("cookie", strings.TrimSpace(strings.Join(keep, ";")))
		}
	}
	c.next.ServeHTTP(rw, req)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func compileRegexes(patterns []string) ([]*regexp.Regexp, error) {
	regexes := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid cookieRegexes pattern %q: %w", pattern, err)
		}
		regexes = append(regexes, reg)
	}

	return regexes, nil
}

func (c *CookieStrip) shouldStrip(cookieName string) bool {
	if stringInSlice(cookieName, c.cookies) {
		return true
	}

	for _, reg := range c.cookieRegexes {
		if reg.MatchString(cookieName) {
			return true
		}
	}

	return false
}
