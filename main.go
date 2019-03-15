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
	var address string
	var requestsNumber int64
	var timeoutMilliseconds float64
	requests := newRequests()
	err := initVariables(&address, &requestsNumber, &timeoutMilliseconds)
	check(err)
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

func initVariables(address *string, requestsNumber *int64, timeout *float64) (err error) {
	addressFlag := flag.String("address", "", "address")
	requestsNumberFlag := flag.String("requestsNumber", "0", "requestsNumber")
	timeoutMillisecondsFlag := flag.String("timeoutMilliseconds", "0", "timeoutMilliseconds")
	flag.Parse()
	*address = *addressFlag
	*requestsNumber, err = strconv.ParseInt(*requestsNumberFlag, 0, 64)
	if err != nil {
		return err
	}
	*timeout, err = strconv.ParseFloat(*timeoutMillisecondsFlag, 64)
	*timeout *= 1000000
	if err != nil {
		return err
	}
	return
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
