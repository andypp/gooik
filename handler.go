package gooik

// Interface for all parsers
type Handler interface {
  Handle(*ContentRequest) (*ContentRequest, <-chan *ContentRequest)
  GenerateNext(*ContentRequest) *ContentRequest
  GenerateChan(*ContentRequest, chan<- *ContentRequest)
}

// Base handler with refactored method
type BaseHandler struct {
  Handler
}

// Provide basic skeleton to handle request
func (h BaseHandler) Handle(con *ContentRequest) (*ContentRequest, <-chan *ContentRequest) {
  contentReqs := make(chan *ContentRequest)

  go func() {
    h.Handler.GenerateChan(con, contentReqs)
    close(contentReqs)
  }()

  return h.Handler.GenerateNext(con), contentReqs
}

// Generate content request
func (h BaseHandler) GenerateNext(con *ContentRequest) *ContentRequest {
  return new(ContentRequest)
}

// Populate content request channel
func (h BaseHandler) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
}
