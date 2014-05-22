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
  expect(t, nextReq, (*ContentRequest)(nil))
  urls := make([]interface{}, 0)

  for contentReq := range contentReqs {
    expect(t, contentReq.Req, (*Request)(nil))
    urls = append(urls, contentReq.Req)
  }
  expect(t, len(urls), 3)

  re = regexp.MustCompile(`<a class="bluelink" href="([^"]+/listing/[^"]+)">`)
  parser = NewLinkParser(re)
  con.Res.Content = `
<a class="bluelink" href="http://www.example.com/listing/abcd">
<a class="bluelink" href="http://www.example.com/listing/efgh/ijkl">
<a class="bluelink" href="http://www.example.com/lis/not-included)">
<a class="bluelink" href="http://www.example.com/listing/12345">
  `

  nextReq, contentReqs = parser.Handle(con)
  expect(t, nextReq, (*ContentRequest)(nil))

  urls = nil

  for contentReq := range contentReqs {
    refute(t, contentReq.Req, (*ContentRequest)(nil))
    urls = append(urls, contentReq.Req.Url.String())
  }
  expect(t, len(urls), 3)


  expect(t, urls[0], "http://www.example.com/listing/abcd")
  expect(t, urls[1], "http://www.example.com/listing/efgh/ijkl")
  expect(t, urls[2], "http://www.example.com/listing/12345")
}
