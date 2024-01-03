package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type METHOD int

func (m *METHOD) String() (error, string) {
	switch *m {
	case GET:
		return nil, "GET"
	case POST:
		return nil, "POST"
	case PATCH:
		return nil, "PATCH"
	case PUT:
		return nil, "PUT"
	case DELETE:
		return nil, "DELETE"
	default:
		return errors.New(fmt.Sprintf("invalid method type `%d`", *m)), "-"
	}
}

const (
	GET METHOD = iota
	POST
	PATCH
	PUT
	DELETE
)

var router *mux.Router
var logger *log.Logger

func Init(aLogger *log.Logger) {
	router = mux.NewRouter()
	logger = aLogger
}
func Listen(port string) {
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}

	logger.Printf("Listening on port `%s`", port)
}

type Handler[T, TE any] func(ctx context.Context, req T) (res TE, err error, status int)

func LoggerMiddleware[T, TE any](
	handler Handler[T, TE],
	logger *log.Logger) (loggedHandler Handler[T, TE]) {
	return func(ctx context.Context, in T) (TE, error, int) {
		start := time.Now()
		logger.Printf("received request with input: `%v`", in)
		res, err, status := handler(ctx, in)
		defer logger.Printf("request processed in `%v` result: `%v`", time.Now().Sub(start), res)

		return res, err, status
	}
}

func AddHandler[T, TE any](
	pattern string,
	handler Handler[T, TE],
	timeout time.Duration,
	methods ...METHOD,
) {
	var m []string
	for _, method := range methods {
		err, mStr := method.String()
		if err != nil {
			logger.Fatal(err)
		}
		m = append(m, mStr)
	}

	router.HandleFunc(pattern, requestHandler[T, TE](LoggerMiddleware[T, TE](handler, logger), timeout, logger)).Methods(m...)
	logger.Printf("registered a new handler to path: `%s` with method: `%s`", pattern, m)
}

func requestHandler[T, TE any](
	handler Handler[T, TE],
	timeout time.Duration,
	logger *log.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var request T

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Printf("error occurred trying to parse the request. Error: `%v`", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err, status := handler(ctx, request)
		if err != nil {
			logger.Printf("returning error response to client. Error: `%v`. Status: `%d`", err, status)
			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Printf("error occurred trying to parse the response. Error: `%v`", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
