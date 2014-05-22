package gooik

import (
  "fmt"
  "net/url"
)

// Request data, nothing but URL at the moment
type Request struct {
  Url *url.URL
}

// Return new instance of Request
func NewRequest(rawurl string) *Request {
  // Make sure URL is valid
  url, err := url.Parse(rawurl)
  if err == nil {
    if url.Scheme != "" && url.Host != "" {
      return &Request{url}
    }
  }
  return nil
}

// Response data
type Response struct {
  Url *url.URL    // URL that was visited
  Content string  // Content of the page
}

// Get base URL of this response URL, with scheme and host only,
// no backslash
func (r Response) BaseUrl() string {
  return  fmt.Sprintf("%v://%v", r.Url.Scheme, r.Url.Host)
}

// Object to contain Request and Response.
// Only used in process content, to indicate if it is necessary to
// download URL before processing content
type ContentRequest struct {
  Req *Request
  Res *Response
}
