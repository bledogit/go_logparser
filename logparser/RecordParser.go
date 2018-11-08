// A record parser
package logparser

import (
	"log"
	"strings"
)

type LogRecord struct {
	path   string
	status string
	ua     string
	cookie string
	ip     string
	good   bool
}

// record parser,  empty record (good = false) on panic
func NewRecord(record string) LogRecord {

	defer func() { // drop funky records
		if rec := recover(); rec != nil {
			log.Println("Recovered from ", record)
		}
	}()

	r := LogRecord{}

	if record[0] == '#' { // skip headers
		r.good = false
		return r
	} else {

		parts := strings.Split(record, "\t")

		r.path = parts[7] + "?" + parts[11]
		r.status = parts[8]
		r.ua = parts[10]
		r.cookie = parts[12]
		r.ip = parts[4]
		r.good = true
		return r
	}
}
