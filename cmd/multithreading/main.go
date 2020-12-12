package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	HTTP()
	Calc()
}

// HTTP example
func HTTP() {
	fmt.Println("== HTTP parallel ==")
	c := make(chan int, 50)
	statusCodes := NewSet()
	start := time.Now()
	for i := 0; i < 50; i++ {
		url := "https://httpbin.org/status/" + strconv.Itoa(200+i)
		go callHTTPAsync(url, c)
	}
	for i := 0; i < 50; i++ {
		statusCodes.Add(<-c)
	}
	duration := time.Since(start)
	fmt.Println("Duration:", duration.Milliseconds())
	if statusCodes.Size() != 50 {
		panic("Not 50 request received")
	}

	fmt.Println()
	fmt.Println("== HTTP serial ==")
	statusCodes2 := NewSet()
	start2 := time.Now()
	for i := 0; i < 50; i++ {
		url := "https://httpbin.org/status/" + strconv.Itoa(200+i)
		statusCodes2.Add(callHTTP(url))
	}
	duration2 := time.Since(start2)
	fmt.Println("Duration:", duration2.Milliseconds())
	if statusCodes2.Size() != 50 {
		panic("Not 50 request received")
	}
}

func callHTTPAsync(url string, c chan (int)) {
	var httpClient = &http.Client{Timeout: time.Second * 1}

	resp, _ := httpClient.Get(url)
	c <- resp.StatusCode
}

func callHTTP(url string) int {
	var httpClient = &http.Client{Timeout: time.Second * 1}

	resp, _ := httpClient.Get(url)
	return resp.StatusCode
}

// Calc example
func Calc() {
	fmt.Println()
	fmt.Println("== Calc serial ==")
	rangeInt := make([]int, 12)
	for i := range rangeInt {
		rangeInt[i] = 30 + i
	}

	var sum = 0
	start := time.Now()
	for _, value := range rangeInt {
		sum += fib(value)
	}
	duration := time.Since(start)
	fmt.Println("Duration:", duration.Milliseconds())

	if sum != 432148168 {
		panic("Sum is wrong")
	}

	fmt.Println()
	fmt.Println("== Calc parallel ==")
	c := make(chan int, 12)
	for _, value := range rangeInt {
		go fibAsync(value, c)
	}

	start2 := time.Now()
	sum = 0
	for range rangeInt {
		sum += <-c
	}
	duration2 := time.Since(start2)
	fmt.Println("Duration:", duration2.Milliseconds())
	if sum != 432148168 {
		panic("Sum is wrong")
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func fibAsync(n int, c chan (int)) {
	if n <= 1 {
		c <- n
	}
	c <- fib(n-1) + fib(n-2)
}

// Set struct
type Set struct {
	list map[int]struct{}
}

// Add int to Set
func (s *Set) Add(v int) {
	s.list[v] = struct{}{}
}

// Size return size of set
func (s *Set) Size() int {
	return len(s.list)
}

// NewSet create new set
func NewSet() *Set {
	s := &Set{}
	s.list = make(map[int]struct{})
	return s
}
