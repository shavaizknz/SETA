package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

const RequestIDKey = "request_id"

func init() {
	// Initialize the logger
	Logger = logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	Logger.SetFormatter(&logrus.JSONFormatter{
		// PrettyPrint: true,
	})
	Logger.SetOutput(os.Stdout)
}

func LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()
		// request Id to be set for request lifecycle
		requestID := generateRequestID()
		c.Set(RequestIDKey, requestID)
		// Create a new context.Context with the request ID
		ctxWithRequestID := context.WithValue(req.Context(), RequestIDKey, requestID)
		// Replace the request context with the one containing the request ID
		c.SetRequest(req.WithContext(ctxWithRequestID))

		// Copy request headers, excluding "Authorization"
		headers := make(http.Header)
		for key, values := range req.Header {
			if key != "Authorization" {
				headers[key] = values
			}
		}

		// Read and capture the request body
		var requestBody interface{}
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)

			// Restore the io.ReadCloser to its original state
			if err == nil {
				c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}

			// Use the captured request body and convert it to interface{}
			if len(bodyBytes) > 0 {
				json.Unmarshal(bodyBytes, &requestBody)
			}
		}

		// obfuscate any sensitive data in the request body (eg. token, email, password)
		requestBody = obfuscate(requestBody)

		// Log request
		requestLog := logrus.Fields{
			"headers":    headers,
			"body":       requestBody,
			"endpoint":   "[" + req.Method + "] " + req.Host + req.URL.String(),
			"request_id": requestID,
			"timestamp":  time.Now().Format(time.RFC3339Nano),
		}

		Logger.WithFields(requestLog).Info("Request Logged")

		//process request
		if err := next(c); err != nil {
			c.Error(err)
		}

		// Log response
		responseLog := logrus.Fields{
			"request_id":  requestID,
			"status_code": res.Status,
			"endpoint":    "[" + req.Method + "] " + req.Host + req.URL.String(),
			"run_time":    time.Since(start).Seconds(),
			"timestamp":   time.Now().Format(time.RFC3339Nano),
		}
		Logger.WithFields(responseLog).Info("Response logged")

		return nil
	}
}

func generateRequestID() string {
	// Generate a UUID as the request ID
	return uuid.New().String()
}

// WithRequestID is a convenience function to create a log entry with request_id
func WithRequestID(ctx interface{}) *logrus.Entry {
	var requestID string

	// Check the type of context
	switch ctx := ctx.(type) {
	case echo.Context:
		// Retrieve request_id from echo.Context
		requestID, _ = ctx.Get(RequestIDKey).(string)
	case context.Context:
		// Retrieve request_id from context.Context
		requestIDValue := ctx.Value(RequestIDKey)
		if requestIDValue != nil {
			requestID, _ = requestIDValue.(string)
		}
	}

	if requestID == "" {
		requestID = "unknown"
	}

	return Logger.WithField("request_id", requestID)
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// obfuscate token, password, email or any sensitive data
func obfuscate(data interface{}) interface{} {
	// Obfuscate sensitive data
	if data == nil {
		return nil
	}

	switch data := data.(type) {
	case map[string]interface{}:
		for key, value := range data {
			switch value := value.(type) {
			case string:
				switch key {
				case "password", "token", "Authorization":
					data[key] = "********"
				case "email":
					data[key] = "********@*****.com"
				default:
					data[key] = value
				}
			default:
				data[key] = obfuscate(value)
			}
		}
	default:
		return data
	}

	return data
}
