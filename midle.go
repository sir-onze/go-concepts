package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// --- given: the core handler ---

func publishHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "published")
}

// --- given: a simple request counter for rate limiting ---
// call counter.Add(1) to increment, counter.Load() to read
var counter atomic.Int64

// --- your job ---

// 1. WithLogging — print method, path, and duration after the request completes
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()           // capture time before
		next.ServeHTTP(w, r)          // run the rest of the chain
		duration := time.Since(start) // calculate how long it took
		fmt.Printf("%s %s took %v\n", r.Method, r.URL.Path, duration)
	})
}

// 2. WithAuth — reject requests with no Authorization header (401)
func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return // ← stop here, next never gets called
		}
		next.ServeHTTP(w, r) // auth passed, continue
	})
}

// 3. WithRateLimit — reject requests once counter exceeds 2 (429)
// hint: increment the counter before checking, then call next or reject
func WithRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := counter.Add(1) // increment and get new value atomically
		if count > 2 {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/publish", WithLogging(WithAuth(WithRateLimit(http.HandlerFunc(publishHandler)))))

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", mux)
}
