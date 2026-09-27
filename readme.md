# Stripcookie
Stripcookie is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which strips cookies by name from a request

### Configuration

### Static

```yaml
pilot:
  token: "xxxxx"

experimental:
  plugins:
    stripcookie:
      moduleName: "github.com/pascal260303/stripcookie-regex"
      version: "v0.1.0"
```

### Dynamic

```yaml
http:
  middlewares:
    strip-foo:
      stripcookie:
        cookies:
          - cookieName
          - otherCookieName
        cookieRegexes:
          - "^session_.*$"
          - "^(utm_|_ga)"
```

`cookies` keeps the existing exact cookie-name matching.

`cookieRegexes` accepts [Go regular expressions](https://pkg.go.dev/regexp/syntax) and each pattern is matched against the cookie name only (not the cookie value).

### AI usage

Parts of this repository were created with AI assistance (GitHub Copilot), with human review and validation before being committed.
