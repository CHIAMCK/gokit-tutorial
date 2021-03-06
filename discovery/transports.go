package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// user will send req to discovery service
// implement client-side discovery
func MakeHttpHandler(endpoint endpoint.Endpoint) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/calculate").Handler(kithttp.NewServer(
		endpoint,
		decodeDiscoverRequest,
		encodeDiscoverResponse,
	))

	return r
}

// ArithmeticRequest define request struct
type ArithmeticRequest struct {
	RequestType string `json:"request_type"`
	A           int    `json:"a"`
	B           int    `json:"b"`
}

// ArithmeticResponse define response struct
type ArithmeticResponse struct {
	Result int   `json:"result"`
	Error  error `json:"error"`
}

func decodeDiscoverRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ArithmeticRequest
	// convert json into struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeDiscoverResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// set header
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	// convert struct into json,
	return json.NewEncoder(w).Encode(response)
}
