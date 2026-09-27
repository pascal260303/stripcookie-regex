package stripcookie_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nilskohrs/stripcookie"
)

func TestDemo(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.Cookies = []string{"testCookie", "otherCookie"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Cookie", "testCookie=testValue; testCookie2=testValue2; test Cookie3=testValue3")
	req.Header.Add("Cookie", "testCookie=testValue; otherCookie=otherValue==abc")

	handler.ServeHTTP(recorder, req)

	assertCookies(t, req, "testCookie2=testValue2; test Cookie3=testValue3")
}

func TestRegexMatching(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.CookieRegexes = []string{`^session_`}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Cookie", "session_id=1; keep=2; session_token=3")

	handler.ServeHTTP(recorder, req)

	assertCookies(t, req, "keep=2")
}

func TestRegexCanMatchMultipleCookieNames(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.CookieRegexes = []string{`^session_.*`}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Cookie", "session_id=1; session_token=2; session_user=3; keep=4")

	handler.ServeHTTP(recorder, req)

	assertCookies(t, req, "keep=4")
}

func TestNonMatchingCookiesArePreserved(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.CookieRegexes = []string{`^session_`}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Cookie", "keepA=1; keepB=2")

	handler.ServeHTTP(recorder, req)

	assertCookies(t, req, "keepA=1; keepB=2")
}

func TestRegexMatchesCookieNameOnly(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.CookieRegexes = []string{`^session_`}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Cookie", "keep=session_id; keep2=abc")

	handler.ServeHTTP(recorder, req)

	assertCookies(t, req, "keep=session_id; keep2=abc")
}

func TestInvalidRegexReturnsError(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.CookieRegexes = []string{"("}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	_, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
	if !strings.Contains(err.Error(), `invalid cookieRegexes pattern "("`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBackwardCompatibilityForExistingConfig(t *testing.T) {
	cfg := stripcookie.CreateConfig()
	cfg.Cookies = []string{"legacyCookie"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := stripcookie.New(ctx, next, cfg, "stripcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Cookie", "legacyCookie=1; keep=2")

	handler.ServeHTTP(recorder, req)

	assertCookies(t, req, "keep=2")
}

func assertCookies(t *testing.T, req *http.Request, expected string) {
	t.Helper()
	if len(req.Header.Values("Cookie")) > 1 {
		t.Errorf("too many headers")
	}
	if req.Header.Get("Cookie") != expected {
		t.Errorf("invalid header value: %s", req.Header.Get("Cookie"))
	}
}
