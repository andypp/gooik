package gooik

import(
  "regexp"
)

// Link parser
type LinkParser struct {
  Regex *regexp.Regexp
}

// Return new link parser
func NewLinkParser(re *regexp.Regexp) *LinkParser {
  return &LinkParser{
    Regex: re,
  }
}

// Parse content, return all links that captured by regex
func (p LinkParser) Handle(con *ContentRequest) (*ContentRequest, <-chan *ContentRequest) {

  // TODO: link parser to build full URL base on request
  contentReqs := make(chan *ContentRequest)
  go func() {
    defer close(contentReqs)
    for _, mt := range p.Regex.FindAllStringSubmatch(con.Res.Content, -1) {
      contentReqs <- &ContentRequest{NewRequest(mt[1]), nil}
    }
  }()
  return nil, contentReqs
}
