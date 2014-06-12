package gooik

import(
  "io/ioutil"
  "net/http"
  "regexp"
  "time"
)

const (
  MAX_CACHE_FILENAME_LENGTH int = 30
)

var (
  nonWordRe *regexp.Regexp = regexp.MustCompile("[^\\w-]")
  dashRe *regexp.Regexp = regexp.MustCompile("--*")
)

// Interface for all downloaders
type Downloader interface {
  Download(req Request) *Response
}

// Download from web
type NetDownloader struct {
}

// Create new NetDownloader instance
func NewNetDownloader() *NetDownloader {
  d := new(NetDownloader)
  return d
}

// Download url from web, with GET method, returns the body as string
func (d *NetDownloader) Download(req Request) *Response {
  // TODO: save download to cache
  response, err := http.Get(req.Url.String())
  if err != nil {
    return nil
  } else {
    defer response.Body.Close()
    //Read and return response body as string
    contents, err := ioutil.ReadAll(response.Body)
    if err == nil {
      if response.StatusCode >= 200 && response.StatusCode < 400 {
        return &Response{req.Url, string(contents)}
      }
    }
    return nil
  }
}

// Mock downloader for testing, read local file
type MockDownloader struct {
  *NetDownloader
  baseDir string
}

// Create instance of MockDownloader
func NewMockDownloader(baseDir string) *MockDownloader {
  md := MockDownloader{
    NetDownloader: NewNetDownloader(),
    baseDir: baseDir,
  }
  return &md
}

// Read file content, sleep to emulate download time
func (d *MockDownloader) Download(req Request) *Response {
  // Simulate download time
  time.Sleep(500 * time.Millisecond)

  cacheFn := GetCacheFileName(req.Url.String())
  dat, err := ioutil.ReadFile(d.baseDir + cacheFn)
  if err != nil {
    return nil
  } else  {
    return &Response{req.Url, string(dat)}
  }
}

// Slug url to save file content
func GetCacheFileName(url string) string {
  if url != "" {
    url = nonWordRe.ReplaceAllString(url, "-")
    url = dashRe.ReplaceAllString(url, "-")
    if len(url) > MAX_CACHE_FILENAME_LENGTH {
      return url[0:MAX_CACHE_FILENAME_LENGTH] + ".html"
    }
    return url + ".html"
  }
  return ""
}
