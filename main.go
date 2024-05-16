package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Payload struct {
	RandNum int
	FibSeq  []int
}

type Metric struct {
	RequestCount   int64
	TotalDuration  float64 // in milliseconds
	AverageLatency float64 // in milliseconds
}

var metrics map[string]*Metric
var PORT, MAX_RANDOM_NUMBER int
var METRIC_DECIMAL_PLACES float64

func main() {
	err := godotenv.Load()
	if err != nil && os.Getenv("PORT") == "" { // if no .env file lookup for environment variables
		log.Fatal("Error loading .env file")
	}

	PORT, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("PORT needs to be numeric: %v", err)
	}
	MAX_RANDOM_NUMBER, err = strconv.Atoi(os.Getenv("APP_MAX_RANDOM_NUMBER"))
	if err != nil {
		log.Fatalf("APP_MAX_RANDOM_NUMBER needs to be numeric: %v", err)
	}
	METRIC_DECIMAL_PLACES, err = strconv.ParseFloat(os.Getenv("APP_METRIC_DECIMAL_PLACES"), 64)
	if err != nil {
		log.Fatalf("APP_METRIC_DECIMAL_PLACES needs to be numeric: %v", err)
	}
	bind := fmt.Sprintf(":%d", PORT)
	// seeding random source
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// metrics map
	metrics = make(map[string]*Metric)
	metrics["health"] = &Metric{0, 0.0, 0.0}
	metrics["ready"] = &Metric{0, 0.0, 0.0}
	metrics["payload"] = &Metric{0, 0.0, 0.0}
	metrics["metrics"] = &Metric{0, 0.0, 0.0}
	// handle functions
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/ready", handleReady)
	http.HandleFunc("/payload", handlePayload)
	http.HandleFunc("/metrics", handleMetrics)
	// start server
	log.Printf("The microservice is listening on %s\n", bind)
	err = http.ListenAndServe(bind, nil)
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
	// update duration
	duration := time.Since(start).Seconds() * 1000 // milliseconds
	m.TotalDuration += duration
	// rounding up to DECIMAL_PLACES
	m.TotalDuration = math.Round(m.TotalDuration*METRIC_DECIMAL_PLACES) / METRIC_DECIMAL_PLACES
	// increase request counter
	m.RequestCount++
	// calculate average latency
	m.AverageLatency = m.TotalDuration / float64(m.RequestCount)
	// rounding up to DECIMAL_PLACES
	m.AverageLatency = math.Round(m.AverageLatency*METRIC_DECIMAL_PLACES) / METRIC_DECIMAL_PLACES
}

func fibonacciDP(n int) []int {
	// dynamic programming version of fibonacci calculation
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
