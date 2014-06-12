package gooik

import (
  "math"
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
  ReturnedUrls map[string]struct{}
  InputUrls []string
}

func (h dummyHandlerOne) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  for _, url := range h.InputUrls {
    h.ReturnedUrls[url] = struct{}{}
  }
}

func newDummyHandlerOne(ret map[string]struct{}, inp []string) dummyHandlerOne {
  d := dummyHandlerOne{
    ReturnedUrls: ret,
    InputUrls: inp,
  }
  d.Handler = d
  return d
}

// Take the first maxElem elements from returned urls
type dummyHandlerTwo struct {
  BaseHandler
  ReturnedUrls map[string]struct{}
  maxElem int
}

func (h dummyHandlerTwo) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  i := 0
  for key := range h.ReturnedUrls {
    if i >= h.maxElem {
      delete(h.ReturnedUrls, key)
    }
    i++
  }
}

func newDummyHandlerTwo(ret map[string]struct{}, max int) dummyHandlerTwo {
  d := dummyHandlerTwo{
    ReturnedUrls: ret,
    maxElem: max,
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
  ReturnedUrls map[string]struct{}
}

func (h dummyHandlerFour) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  h.ReturnedUrls[con.Req.Url.String()] = struct{}{}
}

func newDummyHandlerFour(ret map[string]struct{}) dummyHandlerFour {
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
    "http://www.example.com/5",
    "http://www.example.com/6",
    "http://www.example.com/7",
    "http://www.example.com/8",
    "http://www.example.com/9",
  }
)

func TestSingleListingHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make(map[string]struct{})
  spider.AddListingHandler(
    newDummyHandlerOne(
      urls,
      expectedUrls,
    ),
  )
  spider.Start()

  for _, url := range expectedUrls {
    _, ok := urls[url]
    expect("url found", t, ok, true)
  }
}


func TestDoubleListingHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make(map[string]struct{})
  spider.SetListingHandlers(
    newDummyHandlerOne(
      urls,
      expectedUrls,
    ),
    newDummyHandlerTwo(urls, 3),
  )
  spider.Start()

  expect("url length", t, len(urls), int(math.Min(3, float64(len(expectedUrls)))))
}

func TestTripleListingHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make(map[string]struct{})
  spider.SetListingHandlers(
    newDummyHandlerZero(),
    newDummyHandlerOne(
      urls,
      expectedUrls,
    ),
    newDummyHandlerTwo(urls, 3),
  )
  spider.Start()

  expect("url length", t, len(urls), 0)
}

func TestListingContentHandler(t *testing.T) {
  spider := NewListingSpider("http://www.example.com")
  urls := make(map[string]struct{})
  spider.SetListingHandlers(
    newDummyHandlerThree(expectedUrls),
  )
  spider.AddContentHandler(
    newDummyHandlerFour(urls),
  )
  spider.Start()

  for _, url := range expectedUrls {
    _, ok := urls[url]
    expect("url found", t, ok, true)
  }
}
