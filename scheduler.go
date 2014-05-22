package gooik

import (
  "sync"
)

// Scheduler
type Scheduler struct {
  listingQueue chan Request        // queue for listing pages
  contentQueue chan ContentRequest // queue for content pages
  wait *sync.WaitGroup             // to keep track of both queues
}

// Create new instance of Scheduler
func NewScheduler() *Scheduler {
  s := &Scheduler{
    listingQueue: make(chan Request),
    contentQueue: make(chan ContentRequest),
    wait: new(sync.WaitGroup),
  }
  return s
}

// Generate start URLs
func (s *Scheduler) GenerateStartUrls(urls ...string) {
  for _, url := range urls {
    s.AddToListingQueue(NewRequest(url))
  }
}

// Add url to listing queue
func (s *Scheduler) AddToListingQueue(req *Request) {
  s.addWaiting()
  go func() {
    s.listingQueue <- *req
  }()
}

// Add url to content queue
func (s *Scheduler) AddToContentQueue(con *ContentRequest) {
  s.addWaiting()
  go func() {
    s.contentQueue <- *con
  }()
}

// Increment wait counter
func (s *Scheduler) addWaiting() {
  s.wait.Add(1)
}

// Return listing queue
func (s *Scheduler) ListingQueue() <-chan Request {
  return s.listingQueue
}

// Return content queue
func (s *Scheduler) ContentQueue() <-chan ContentRequest {
  return s.contentQueue
}

// Wait until all queues are empty
func (s *Scheduler) Wait() {
  s.wait.Wait()
}

// Decrement waiting counter
func (s *Scheduler) Done() {
  s.wait.Done()
}
