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
  refute(t, req, (*Request)(nil))

  req = NewRequest(fmt.Sprintf("%v://%v%v?%v#%v", scheme, host, path, query, fragment))
  refute(t, req, (*Request)(nil))
  expect(t, req.Url.Scheme, scheme)
  expect(t, req.Url.Host, host)
  expect(t, req.Url.Path, path)
  expect(t, req.Url.RawQuery, query)
  expect(t, req.Url.Fragment, fragment)

  req = NewRequest("this-is-invalid-url")
  expect(t, req, (*Request)(nil))
}

func TestResponseBaseUrl(t *testing.T) {
  req := NewRequest("http://www.example.com/en/search?walk=123")

  res := Response{
    Url: req.Url,
    Content: "this is some content string",
  }
  expect(t, res.BaseUrl(), "http://www.example.com")
}
