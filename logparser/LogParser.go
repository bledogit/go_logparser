//
// LogParser, parses cloudfront logs from S3 log file,
// reconstruct the original get request, recovering the  original ip, cookie and user agent
// resending it to a new endpoint (URL)
//
// s3 -> +--> worker -->+-->result
//       +--> worker -->+
//       ...            |
//       +--> worker -->+
//
// Records in s3 file are distributed across a worker pool
//
package logparser

import (
	"bufio"
	"log"
	"strconv"
	"sync"

	"compress/gzip"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ( // global s3 service
	svc  *s3.S3
	once sync.Once
)

type LogParser struct {
	maxRequests int
	base        string
	records     chan LogRecord
	results     chan result
	nWorkers    int
	stats       map[string]int
}

type result struct {
	code int
}

// factory
func s3Service() *s3.S3 {
	once.Do(func() {

		svc = s3.New(session.New(&aws.Config{
			Region: aws.String("us-east-1")}))
	})
	return svc
}

// New parser
// defaults to 1000 workers, all requests
// @param base Base url of endpoint
func NewParser(base string) LogParser {
	n := 1000

	l := LogParser{}
	l.maxRequests = 999999999
	l.base = base
	l.records = make(chan LogRecord, n)
	l.results = make(chan result, n)
	l.nWorkers = n
	l.stats = make(map[string]int)

	return l
}

// Configures parser with a maximum number of requests
func (l *LogParser) WithMaxRequest(max int) *LogParser {
	l.maxRequests = max
	return l
}

// Configures parser with a maximum number of workers
func (l *LogParser) WithMaxWorkers(max int) *LogParser {
	l.records = make(chan LogRecord, max)
	l.results = make(chan result, max)
	l.nWorkers = max
	return l
}

// Get statistics (map with error code, number of entries that resulted in this condition)
func (l *LogParser) GetStats() map[string]int {
	return l.stats
}

// Parses a s3 Object
func (l *LogParser) ParseS3Object(entity events.S3Entity) {
	l.ParseS3ObjectKey(entity.Object.Key, entity.Bucket.Name)
}

// Parses a s3 Object given bucket and key
func (l *LogParser) ParseS3ObjectKey(object string, bucket string) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	}

	result, err := s3Service().GetObject(input)
	if err != nil {
		log.Fatal(err)
	}
	unzip, err := gzip.NewReader(result.Body)
	if err != nil {
		log.Fatal(err)
	}
	liner := bufio.NewScanner(unzip)

	var wg sync.WaitGroup
	wg.Add(1)
	go l.createWorkerPool()
	go l.processResuts(&wg)

	count := 0

	for liner.Scan() {
		count++
		r := NewRecord(liner.Text())
		l.records <- r
		if count >= l.maxRequests {
			break
		}
	}

	close(l.records)

	if err := liner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

}

// picks records from record channel into results channel which contains the http error code
func (l *LogParser) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	sender := NewSender()
	count := 0
	for record := range l.records {
		count++
		if record.good {
			l.results <- result{sender.sendRecord(l.base, record)}
		} else {
			l.results <- result{-1}
		}
	}

	//log.Println(id, " count = ", count)

}

// Creates worker pool and waits for workers to be done
func (l *LogParser) createWorkerPool() {
	//log.Println("Create ", l.nWorkers, " workers")
	var wg sync.WaitGroup
	for i := 0; i < l.nWorkers; i++ {
		wg.Add(1)
		go l.worker(i, &wg)
	}
	wg.Wait()
	close(l.results)
}

// gets all results to mark the end of the processing
func (l *LogParser) processResuts(wg *sync.WaitGroup) {

	defer wg.Done()
	total := 0
	for result := range l.results {
		total++
		code := strconv.Itoa(result.code)
		s, ok := l.stats[code]
		if ok {
			l.stats[code] = s + 1
		} else {
			l.stats[code] = 1
		}
		//if total%100 == 0 {
		//	log.Println(total)
		//}
	}
	l.stats["total"] = total

}
