package gooik

import (
  "testing"
)

// Doesn't do anything, will stop pipeline execution
type dummyHandlerZero struct {
  BaseHandler
}

func (h dummyHandlerZero) GenerateNext(con *ContentRequest) *ContentRequest {
  return nil
}

func newDummyHandlerZero() dummyHandlerZero {
  d := dummyHandlerZero{}
  d.Handler = d
  return d
}

// Copy urls from input urls to returned urls
type dummyHandlerOne struct {
  BaseHandler
  ReturnedUrls *[]string
  InputUrls []string
}

func (h dummyHandlerOne) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  for _, url := range h.InputUrls {
    *h.ReturnedUrls = append(*h.ReturnedUrls, url)
  }
}

func newDummyHandlerOne(ret *[]string, inp []string) dummyHandlerOne {
  d := dummyHandlerOne{
    ReturnedUrls: ret,
    InputUrls: inp,
  }
  d.Handler = d
  return d
}

// Take the first three elements from returned urls
type dummyHandlerTwo struct {
  BaseHandler
  ReturnedUrls *[]string
}

func (h dummyHandlerTwo) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  *h.ReturnedUrls = (*h.ReturnedUrls)[0:3]
}

func newDummyHandlerTwo(ret *[]string) dummyHandlerTwo {
  d := dummyHandlerTwo{
    ReturnedUrls: ret,
  }
  d.Handler = d
  return d
}

// Pass input urls into content request channel
type dummyHandlerThree struct {
  BaseHandler
  InputUrls []string
}

func (h dummyHandlerThree) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  for _, url := range h.InputUrls {
    channel <- &ContentRequest{
                 Req: NewRequest(url),
               }
  }
}

func newDummyHandlerThree(inp []string) dummyHandlerThree {
  d := dummyHandlerThree{
    InputUrls: inp,
  }
  d.Handler = d
  return d
}

// Append content requests to returned urls
type dummyHandlerFour struct {
  BaseHandler
  ReturnedUrls *[]string
}

func (h dummyHandlerFour) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  *h.ReturnedUrls = append(*h.ReturnedUrls, con.Req.Url.String())
}

func newDummyHandlerFour(ret *[]string) dummyHandlerFour {
  d := dummyHandlerFour{
    ReturnedUrls: ret,
  }
  d.Handler = d
  return d
}

var (
  expectedUrls = []string{
    "http://www.example.com/1",
    "http://www.example.com/2",
    "http://www.example.com/3",
    "http://www.example.com/4",
  }
)

func TestSingleListingHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make([]string, 0)
  spider.AddListingHandler(
    newDummyHandlerOne(
      &urls,
      expectedUrls,
    ),
  )
  spider.Start()

  expect("url length", t, len(urls), 4)
}


func TestDoubleListingHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make([]string, 0)
  spider.SetListingHandlers(
    newDummyHandlerOne(
      &urls,
      expectedUrls,
    ),
    newDummyHandlerTwo(&urls),
  )
  spider.Start()

  expect("url length", t, len(urls), 3)
}

func TestTripleListingHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make([]string, 0)
  spider.SetListingHandlers(
    newDummyHandlerZero(),
    newDummyHandlerOne(
      &urls,
      expectedUrls,
    ),
    newDummyHandlerTwo(&urls),
  )
  spider.Start()

  expect("url length", t, len(urls), 0)
}

func TestListingContentHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make([]string, 0)
  spider.SetListingHandlers(
    newDummyHandlerThree(expectedUrls),
  )
  spider.AddContentHandler(
    newDummyHandlerFour(&urls),
  )
  spider.Start()

  expect("url length", t, len(urls), 4)
}
