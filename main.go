package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	requests := newRequests()
	var wg sync.WaitGroup
	for i := int64(0); i < requestsNumber; i++ {
		wg.Add(1)
		go func() {
			requestStart := time.Now()
			client := http.Client{Timeout: time.Duration(timeoutMilliseconds)}
			_, err := client.Get(address)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				incRejectedNumber(requests)
			} else if err != nil {
				panic(err)
			} else {
				addRequestTime(requests, requestStart)
			}
			wg.Done()
		}()
		wg.Wait()
	}
	printResult(requests)
}

var (
	address             string
	requestsNumber      int64
	timeoutMilliseconds float64
)

const (
	defaultAddress             = "https://www.google.com/"
	defaultRequestsNumber      = 10
	defaultTimeoutMilliseconds = 200
)

func init() {
	address = *flag.String("address", defaultAddress, "address")
	requestsNumberFlag := flag.String("requestsNumber", string(defaultRequestsNumber), "requestsNumber")
	timeoutMillisecondsFlag := flag.String("timeoutMilliseconds", string(defaultTimeoutMilliseconds), "timeoutMilliseconds")
	flag.Parse()
	var err error
	if address == "" {
		address = defaultAddress
		fmt.Printf("address default value is:%s , because address has incorrect value\n", address)
	}
	requestsNumber, err = strconv.ParseInt(*requestsNumberFlag, 0, 64)
	if err != nil || requestsNumber <= 0 {
		requestsNumber = defaultRequestsNumber
		fmt.Printf("requestsNumber default value is:%d , because requestsNumber has incorrect value\n", requestsNumber)
	}
	timeoutMilliseconds, err = strconv.ParseFloat(*timeoutMillisecondsFlag, 64)
	timeoutMilliseconds *= 1000000
	if err != nil || requestsNumber <= 0 {
		timeoutMilliseconds = defaultTimeoutMilliseconds
		fmt.Printf("timeoutMilliseconds default value is: %f, because timeoutMilliseconds has incorrect value\n", timeoutMilliseconds)
	}
}

type request struct {
	mux            sync.Mutex
	requestTimes   []time.Duration
	numberRejected int64
}

func newRequests() *request {
	return &request{requestTimes: []time.Duration{}}
}

func maxTime(requestsTime []time.Duration) (maxTime time.Duration) {
	if len(requestsTime) <= 0 {
		return
	}
	maxTime = requestsTime[0]
	for i := 0; i < len(requestsTime); i++ {
		if requestsTime[i] > maxTime {
			maxTime = requestsTime[i]
		}
	}
	return
}
func minTime(requestsTime []time.Duration) (minTime time.Duration) {
	if len(requestsTime) <= 0 {
		return
	}
	minTime = requestsTime[0]
	for i := 0; i < len(requestsTime); i++ {
		if requestsTime[i] < minTime {
			minTime = requestsTime[i]
		}
	}
	return
}
func requestsAverageTime(requestTimes []time.Duration) time.Duration {
	if len(requestTimes) != 0 {
		return sum(requestTimes) / time.Duration(len(requestTimes))
	}
	return 0
}
func sum(times []time.Duration) (sum time.Duration) {
	for i := 0; i < len(times); i++ {
		sum += times[i]
	}
	return
}
func incRejectedNumber(requests *request) {
	requests.mux.Lock()
	requests.numberRejected++
	requests.mux.Unlock()
}
func addRequestTime(requests *request, requestStart time.Time) {
	requests.mux.Lock()
	requests.requestTimes = append(requests.requestTimes, time.Since(requestStart))
	requests.mux.Unlock()
}
func printResult(requests *request) {
	fmt.Println("End time of requests:", sum(requests.requestTimes))
	fmt.Println("Average request time:", requestsAverageTime(requests.requestTimes))
	fmt.Println("Longest request time:", maxTime(requests.requestTimes))
	fmt.Println("Faster request time:", minTime(requests.requestTimes))
	fmt.Println("Responds number that didn't wait:", requests.numberRejected)
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
