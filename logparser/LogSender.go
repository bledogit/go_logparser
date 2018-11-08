// LogSender
// Manages a connection pool with trasnport set to keep connection alive
//
package logparser

import (
	"io/ioutil"
	"log"
	"net/http"
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
func (s *Sender) sendRecord(baseUrl string, record LogRecord) int {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from ", r)
		}
	}()

	url := baseUrl + record.path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Can't request")
	}

	req.Header.Add("X-Forwarded-For", record.ip)
	req.Header.Add("User-Agent", record.ua)
	req.Header.Add("X-Device-User-Agent", record.ua)
	req.Header.Add("Cookie", record.cookie)
	//req.Close = true

	resp2, err := s.client.Do(req)
	if err != nil {
		panic("Can't do shit" + err.Error())
	}
	defer resp2.Body.Close()

	// don't care about the response
	ioutil.ReadAll(resp2.Body)

	return resp2.StatusCode

}
