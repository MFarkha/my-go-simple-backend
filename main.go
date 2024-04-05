package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	MAX_RANDOM_NUMBER = 20
	PORT              = 3000
)

type Payload struct {
	RandNum int
	FibSeq  []int
}

type Metric struct {
	RequestCount   int
	TotalDuration  float64
	AverageLatency float64
}

var metrics map[string]*Metric

func main() {
	bind := fmt.Sprintf(":%d", PORT)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	metrics = make(map[string]*Metric)
	metrics["health"] = &Metric{0, 0.0, 0.0}
	metrics["ready"] = &Metric{0, 0.0, 0.0}
	metrics["payload"] = &Metric{0, 0.0, 0.0}
	metrics["metrics"] = &Metric{0, 0.0, 0.0}

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/ready", handleReady)
	http.HandleFunc("/payload", handlePayload)
	http.HandleFunc("/metrics", handleMetrics)
	log.Printf("The microservice version 2 is listening on %s\n", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatalf("The microservice failure to bind: %v, error: %v", bind, err)
	}
}

func handleHealth(w http.ResponseWriter, req *http.Request) {
	// updating metrics for the endpoint
	defer updateMetric(metrics["health"], time.Now())
	// sending request
	returnJson(w, func() (interface{}, error) {
		return "Service is healthy", nil
	})
}

func handleReady(w http.ResponseWriter, req *http.Request) {
	// updating metrics for the endpoint
	defer updateMetric(metrics["ready"], time.Now())
	// sending request
	returnJson(w, func() (interface{}, error) {
		return "Service is ready", nil
	})
}

func handlePayload(w http.ResponseWriter, req *http.Request) {
	// updating metrics for the endpoint
	defer updateMetric(metrics["payload"], time.Now())
	// calculating fibonacci sequence
	n := rand.Intn(MAX_RANDOM_NUMBER)
	fibSeq := fibonacciDP(n)
	// response
	returnJson(w, func() (interface{}, error) {
		return &Payload{
			RandNum: n,
			FibSeq:  fibSeq,
		}, nil
	})
}

func handleMetrics(w http.ResponseWriter, req *http.Request) {
	// updating metrics for the endpoint
	defer updateMetric(metrics["metrics"], time.Now())
	// Send metrics response
	returnJson(w, func() (interface{}, error) {
		return metrics, nil
	})
}

func updateMetric(m *Metric, start time.Time) {
	duration := time.Since(start).Seconds()
	m.TotalDuration += duration
	m.RequestCount++
	m.AverageLatency = m.TotalDuration / float64(m.RequestCount)
}

func fibonacciDP(n int) []int {
	seq := make([]int, n+2)
	seq[0], seq[1] = 0, 1
	for i := 2; i <= n; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	return seq[:n+1]
}

func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, serverErr := withData()
	if serverErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Printf("error from json.marshal of serverErr: %v", err)
			return err
		}
		w.Write(serverErrJson)
		return nil
	}
	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Printf("error from json.marshal of data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Write(dataJson)
	return nil
}
