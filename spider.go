package gooik

// Base spider
type Spider struct {
  scheduler *Scheduler
  listingHandlers []Handler
  contentHandlers []Handler
}

// Return new instance of Spider
func NewSpider() *Spider {
  s := &Spider{
    scheduler: NewScheduler(),
  }
  return s
}

// Add listing handler
func (s *Spider) AddListingHandler(h Handler) {
  s.listingHandlers = append(s.listingHandlers, h)
}

// Clear listing handlers and add new ones
func (s *Spider) SetListingHandlers(hs...Handler) {
  s.listingHandlers = make([]Handler, 0)
  for _, h := range hs {
    s.AddListingHandler(h)
  }
}

// Add content handler
func (s *Spider) AddContentHandler(h Handler) {
  s.contentHandlers = append(s.contentHandlers, h)
}

// Clear content handlers and add new ones
func (s *Spider) SetContentHandlers(hs...Handler) {
  s.contentHandlers = make([]Handler, 0)
  for _, h := range hs {
    s.AddContentHandler(h)
  }
}

// Listing spider
type ListingSpider struct {
  *Spider
}

// Return new instance of ListingSpider
func NewListingSpider(startPage string) *ListingSpider {
  ls := ListingSpider{
    Spider: NewSpider(),
  }
  ls.scheduler.GenerateStartUrls(startPage)
  return &ls
}

// Start crawling until all queues are empty
func (s *Spider) Start() {
  // Process all items in listing queue
  go func() {
    for req := range s.scheduler.ListingQueue() {
      go s.handleListing(req)
    }
  }()

  // Process all items in content queue
  go func() {
    for con := range s.scheduler.ContentQueue() {
      // Copy request, and pass pointer to the copy
      // this is to avoid all goroutines in this iteration
      // processing the same request
      temp := con
      go s.handleContent(&temp)
    }
  }()

  // Wait until both queues are empty
  s.scheduler.Wait()
}

// Handle listing item with all handlers
func (s *Spider) handleListing(req Request) {
  // Set the request for first handler
  nowReq := &ContentRequest{
    Req: &req,
  }
  for _, handler := range s.listingHandlers {
    nextReq := s.doHandle(handler, nowReq)
    if nextReq == nil {
      break
    }
    // Pass this request to next handler
    nowReq = nextReq
  }
  s.scheduler.Done()
}

// Handle content item with all handlers
func (s *Spider) handleContent(nowReq *ContentRequest) {
  for _, handler := range s.contentHandlers {
    nextReq := s.doHandle(handler, nowReq)
    if nextReq == nil {
      break
    }
    // Pass this request to next handler
    nowReq = nextReq
  }
  s.scheduler.Done()
}

// Perform actual content handling for each handler
func (s *Spider) doHandle(handler Handler, nowReq *ContentRequest) *ContentRequest {
  nextReq, contentReqs := handler.Handle(nowReq)

  // Queue content requests
  for contentReq := range contentReqs {
    s.scheduler.AddToContentQueue(contentReq)
  }

  if nextReq != nil {
    // Queue next listing page
    if nextReq.Req != nil {
      s.scheduler.AddToListingQueue(nextReq.Req)
    }
    return nextReq
  }
  return nil
}
