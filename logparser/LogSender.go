// LogSender
// Manages a connection pool with trasnport set to keep connection alive
//
package logparser

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Sender struct {
	client *http.Client
}

func NewSender() *Sender {

	s := Sender{}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}

	s.client = &http.Client{Transport: tr}

	return &s
}

// send a record, 0 on panic
func (s *Sender) sendRecord(baseUrl string, record LogRecord) Result {

	result := Result{0, "good"}

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from sendRecord", r)
			result.err = "unknown"
			result.code = -3
		}
	}()

	targetUrl := baseUrl + record.path

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err == nil {
		req.Header.Add("X-Forwarded-For", record.ip)
		req.Header.Add("User-Agent", record.ua)
		req.Header.Add("X-Device-User-Agent", record.ua)
		req.Header.Add("Cookie", record.cookie)
		//req.Close = true
		resp, err := s.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			// don't care about the response
			ioutil.ReadAll(resp.Body)
			result.code = resp.StatusCode
		} else {
			if serr, ok := err.(*url.Error); ok {
				result.err = serr.Err.Error()
			} else {
				result.err = "unknown(req)"
			}
			result.code = -2
		}
	} else {
		if serr, ok := err.(*url.Error); ok {
			result.err = serr.Err.Error()
		} else {
			result.err = "unknown(pool)"
		}
		result.code = -1
	}

	return result

}
