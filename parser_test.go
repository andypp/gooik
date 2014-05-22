package gooik

import (
  "regexp"
  "testing"
)

func TestLinkParser(t *testing.T) {
  re := regexp.MustCompile(`<a class="bluelink" href="/listing/([^"]+)">`)
  parser := NewLinkParser(re)
  req := NewRequest("http://www.example.com")
  con := &ContentRequest{
    Req: req,
    Res: &Response{
           Url: req.Url,
           Content: `
<a class="bluelink" href="/listing/abcd">
<a class="bluelink" href="/listing/efgh/ijkl">
<a class="bluelink" href="/lis/not-included)">
<a class="bluelink" href="/listing/12345">
           `,
         },
  }

  nextReq, contentReqs := parser.Handle(con)
  expect("next request 1", t, nextReq, (*ContentRequest)(nil))
  urls := make([]interface{}, 0)

  for contentReq := range contentReqs {
    expect("content request 1", t, contentReq.Req, (*Request)(nil))
    urls = append(urls, contentReq.Req)
  }
  expect("url length 1", t, len(urls), 3)

  re = regexp.MustCompile(`<a class="bluelink" href="([^"]+/listing/[^"]+)">`)
  parser = NewLinkParser(re)
  con.Res.Content = `
<a class="bluelink" href="http://www.example.com/listing/abcd">
<a class="bluelink" href="http://www.example.com/listing/efgh/ijkl">
<a class="bluelink" href="http://www.example.com/lis/not-included)">
<a class="bluelink" href="http://www.example.com/listing/12345">
  `

  nextReq, contentReqs = parser.Handle(con)
  expect("next request 2", t, nextReq, (*ContentRequest)(nil))

  urls = nil

  for contentReq := range contentReqs {
    refute("content request 2", t, contentReq.Req, (*ContentRequest)(nil))
    urls = append(urls, contentReq.Req.Url.String())
  }
  expect("url length 2", t, len(urls), 3)


  expect("url 1", t, urls[0], "http://www.example.com/listing/abcd")
  expect("url 2", t, urls[1], "http://www.example.com/listing/efgh/ijkl")
  expect("url 3", t, urls[2], "http://www.example.com/listing/12345")
}
