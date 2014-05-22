package gooik

import (
  "fmt"
  "testing"
)

func retrieveUrls(urls *[]string, scheduler *Scheduler, ty string) {
  if ty == "listing" {
    for req := range scheduler.ListingQueue() {
      *urls = append(*urls, req.Url.String())
      scheduler.Done()
    }
  } else if ty == "content" {
    for con := range scheduler.ContentQueue() {
      *urls = append(*urls, con.Req.Url.String())
      scheduler.Done()
    }
  } else {
    panic(fmt.Sprintf("Invalid type: %s", ty))
  }
}

func TestGenerateStartUrls(t *testing.T) {
  scheduler := NewScheduler()
  scheduler.GenerateStartUrls(
    "http://www.example.com/1",
    "http://www.example.com/2",
    "http://www.example.com/3",
  )

  urls := make([]string, 0)
  go retrieveUrls(&urls, scheduler, "listing")
  scheduler.Wait()
  expect("url 1", t, urls[0], "http://www.example.com/1")
  expect("url 2", t, urls[1], "http://www.example.com/2")
  expect("url 3", t, urls[2], "http://www.example.com/3")
}

func TestAddListingQueue(t *testing.T) {
  scheduler := NewScheduler()
  expectedUrls := []string{
    "http://www.example.com/1",
    "http://www.example.com/2",
    "http://www.example.com/3",
  }
  for _, eu := range expectedUrls {
    scheduler.AddToListingQueue(NewRequest(eu))
  }

  urls := make([]string, 0)
  go retrieveUrls(&urls, scheduler, "listing")
  scheduler.Wait()
  expect("url length", t, len(urls), len(expectedUrls))
  for i := 0; i < len(expectedUrls); i++ {
    expect("url", t, urls[i], expectedUrls[i])
  }
}

func TestAddContentQueue(t *testing.T) {
  scheduler := NewScheduler()
  expectedUrls := []string{
    "http://www.example.com/1",
    "http://www.example.com/2",
    "http://www.example.com/3",
  }
  for _, eu := range expectedUrls {
    scheduler.AddToContentQueue(&ContentRequest{NewRequest(eu), nil})
  }

  urls := make([]string, 0)
  go retrieveUrls(&urls, scheduler, "content")
  scheduler.Wait()
  expect("url length", t, len(urls), len(expectedUrls))
  for i := 0; i < len(expectedUrls); i++ {
    expect("url", t, urls[i], expectedUrls[i])
  }
}
