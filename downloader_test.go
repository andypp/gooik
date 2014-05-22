package gooik

import (
  "fmt"
	"net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

type DummyHttpHandler struct{
  bodyString string
}

func (e DummyHttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	sts := r.FormValue("sts")

	if sts == "404" {
		rw.WriteHeader(404)
	} else {
		rw.Write([]byte(e.bodyString))
	}
}

func TestNetDownloader(t *testing.T) {
  bodyString := "this is a body"
  server := httptest.NewServer(&DummyHttpHandler{bodyString})
  defer server.Close()

  downloader := NewNetDownloader()
  reqUrl := fmt.Sprintf("%s?walk=123", server.URL)
  req := NewRequest(reqUrl)

  res := downloader.Download(*req)
  refute(t, res, (*Response)(nil))
  expect(t, res.Url.String(), reqUrl)
  expect(t, res.Content, bodyString)

  req404 := NewRequest(reqUrl + "&sts=404")

  res404 := downloader.Download(*req404)
  expect(t, res404, (*Response)(nil))
}

func TestCacheFileName(t *testing.T) {
  cacheFn := GetCacheFileName("http://www.example.com")
  expect(t, cacheFn, "http-www-example-com.html")

  cacheFn = GetCacheFileName("http://www.example.com?walk=tes++++12")
  expect(t, cacheFn, "http-www-example-com-walk-tes-.html")

  cacheFn = GetCacheFileName("http://www.gle.com?walk=test+.....+++++12")
  expect(t, cacheFn, "http-www-gle-com-walk-test-12.html")
}

func TestMockDownloader(t *testing.T) {
  downloader := NewMockDownloader("")

  rawUrl := "http://testing.cache?file"
  req := NewRequest(rawUrl)

  res := downloader.Download(*req)
  refute(t, res, (*Response)(nil))
  expect(t, res.Url.String(), rawUrl)

  trimmedContent := strings.Trim(res.Content, " \n")
  expect(t, trimmedContent, "This is a dummy content")

  rawUrl = "http://testing.cache?file21"
  req = NewRequest(rawUrl)
  res = downloader.Download(*req)
  expect(t, res, (*Response)(nil))
}