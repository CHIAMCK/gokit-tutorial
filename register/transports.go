package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// using golang built-in library to create error
// shorthand for defining variables in bulk
// define global variables that can be used by all the functions in this package
var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// expose your service to the outside world
// create http handler using mux
func MakeHttpHandler(ctx context.Context, endpoints ArithmeticEndpoints, logger log.Logger) http.Handler {
	// create router instance
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		// to encode errors to the http.ResponseWriter
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}

	// mux Handler accepts struct that implements ServeHTTP
	r.Methods("POST").Path("/calculate/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoints.ArithmeticEndpoint,
		decodeArithmeticRequest,
		encodeArithmeticResponse,
		options...,
	))

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeArithmeticResponse,
		options...,
	))

	return r
}

// func MakeHttpHandler(endpoint endpoint.Endpoint) http.Handler {
// 	r := mux.NewRouter()

// 	r.Methods("POST").Path("/calculate").Handler(kithttp.NewServer(
// 		endpoint,
// 		decodeDiscoverRequest,
// 		encodeDiscoverResponse,
// 	))

// 	return r
// }

// create a function to convert payload to struct
func decodeArithmeticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// access URL params
	vars := mux.Vars(r)
	// second return value indicates whether the key is found
	requestType, ok := vars["type"]
	// if type param is not passed
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	// validation

	pb, ok := vars["b"]
	// validation

	// strconv stands for string conversion
	// convert string to int
	a, _ := strconv.Atoi(pa)
	b, _ := strconv.Atoi(pb)

	// return the struct
	return ArithmeticRequest{
		RequestType: requestType,
		A:           a,
		B:           b,
	}, nil
}

// generate response
// convert struct into json
func encodeArithmeticResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// set header
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	// encode and write data to the underlying output steam in a single step
	return json.NewEncoder(w).Encode(response)
}

func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}

func decodeDiscoverRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ArithmeticRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeDiscoverResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
