package main

import (
	"context"
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type ArithmeticEndpoints struct {
	ArithmeticEndpoint  endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

var (
	ErrInvalidRequestType = errors.New("RequestType has only four type: Add,Subtract,Multiply,Divide")
)

type ArithmeticRequest struct {
	RequestType string `json:"request_type"`
	A           int    `json:"a"`
	B           int    `json:"b"`
}

type ArithmeticResponse struct {
	Result int   `json:"result"`
	Error  error `json:"error"`
}

// type ArithmeticEndpoint endpoint.Endpoint

// endpoint represents a single RPC method.
// call the service functions here to process the request
// an adapter to convert each of our service's methods into an endpoint
// each adapter takes a service method and returns an endpoint
func MakeArithmeticEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// type assertion
		// provides access to an interface value's underlying concrete value
		// interface value holds the concrete type ArithmeticRequest
		req := request.(ArithmeticRequest)

		var (
			res, a, b int
			calError  error
		)

		a = req.A
		b = req.B

		// check whether these 2 strings are equal
		if strings.EqualFold(req.RequestType, "Add") {
			res = svc.Add(a, b)
		} else if strings.EqualFold(req.RequestType, "Subtract") {
			res = svc.Subtract(a, b)
		} else if strings.EqualFold(req.RequestType, "Multiply") {
			res = svc.Multiply(a, b)
		} else if strings.EqualFold(req.RequestType, "Divide") {
			res, calError = svc.Divide(a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return ArithmeticResponse{Result: res, Error: calError}, nil
	}
}

type HealthRequest struct{}

type HealthResponse struct {
	Status bool `json:"status"`
}

// accept service
// convert service into endpoint
func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// call service function
		status := svc.HealthCheck()
		// return response
		return HealthResponse{status}, nil
	}
}
