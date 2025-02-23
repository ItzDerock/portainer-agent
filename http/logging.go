package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Create log file
var logFile, _ = os.OpenFile("./agentlogs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
var logger = log.New(logFile, "", log.LstdFlags)

// LoggingMiddleware logs each request and response
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		ip := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		headers := formatHeaders(r.Header)

		// Read and log the request body
		var requestBody bytes.Buffer
		if r.Body != nil {
			_, err := io.Copy(&requestBody, r.Body)
			if err != nil {
				logger.Printf("Error reading request body: %v", err)
			}
			// Restore the body so it can be read again by the handler
			r.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))
		}
		body := requestBody.String()

		logger.Printf("Request from %s\n%s %s\nHeaders:\n%s\nBody:\n%s\n\n", ip, method, path, headers, body)

		// Response logging
		rec := &ResponseRecorder{ResponseWriter: w, body: new(bytes.Buffer), statusCode: http.StatusOK}

		start := time.Now()
		next.ServeHTTP(rec, r)
		duration := time.Since(start)

		// Log response details
		responseHeaders := formatHeaders(rec.Header())
		responseBody := rec.body.String()

		logger.Printf("Response: %d %s\nHeaders:\n%s\nBody:\n%s\n\nDuration: %v\n---\n",
			rec.statusCode, http.StatusText(rec.statusCode), responseHeaders, responseBody, duration)
	})
}

// ResponseRecorder captures response details
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *ResponseRecorder) Write(data []byte) (int, error) {
	// Capture the response body
	rec.body.Write(data)
	return rec.ResponseWriter.Write(data)
}

// formatHeaders formats headers as a string
func formatHeaders(headers http.Header) string {
	var buffer bytes.Buffer
	for key, values := range headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, values))
	}
	return buffer.String()
}
