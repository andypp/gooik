package gooik

import(
  "regexp"
)

// Link parser
type LinkParser struct {
  BaseHandler
  Regex *regexp.Regexp
}

// Return new link parser
func NewLinkParser(re *regexp.Regexp) LinkParser {
  p := LinkParser{
    Regex: re,
  }
  p.Handler = p
  return p
}

// Parse content, return all links that captured by regex
func (p LinkParser) GenerateChan(con *ContentRequest, channel chan<- *ContentRequest) {
  for _, mt := range p.Regex.FindAllStringSubmatch(con.Res.Content, -1) {
    channel <- &ContentRequest{NewRequest(mt[1]), nil}
  }
}

// Return nil content request
func (p LinkParser) GenerateNext(con *ContentRequest) *ContentRequest {
  return nil
}
