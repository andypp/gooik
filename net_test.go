package gooik

import (
  "fmt"
  "testing"
)

func TestNewRequest(t *testing.T) {
  scheme := "http"
  host := "www.example.com"
  path := "/en/search"
  query := "walk=test&go=now"
  fragment := "123"

  req := NewRequest(fmt.Sprintf("%v://%v", scheme, host))
  refute("request 1", t, req, (*Request)(nil))

  req = NewRequest(fmt.Sprintf("%v://%v%v?%v#%v", scheme, host, path, query, fragment))
  expect("scheme", t, req.Url.Scheme, scheme)
  expect("host", t, req.Url.Host, host)
  expect("path", t, req.Url.Path, path)
  expect("query", t, req.Url.RawQuery, query)
  expect("fragment", t, req.Url.Fragment, fragment)

  req = NewRequest("this-is-invalid-url")
  expect("request 2", t, req, (*Request)(nil))
}

func TestResponseBaseUrl(t *testing.T) {
  req := NewRequest("http://www.example.com/en/search?walk=123")

  res := Response{
    Url: req.Url,
    Content: "this is some content string",
  }
  expect("base url", t, res.BaseUrl(), "http://www.example.com")
}
