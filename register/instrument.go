package main

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("Rate limit exceed!")

// Middleware is a chainable behavior modifier for endpoints.
func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	// middleware is a custom function that accepts endpoint and return endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		// endpoint is a custom function that accepts context and request
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// add rate limit checking before process the request
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			// call the next endpoint
			return next(ctx, request)
		}
	}
}

// create a function to return middleware
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	// Middleware is a custom function type that accepts endpoint and return endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// add rate limit checking before process the request
			// check whether n events may happen at time now
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			// proceed to process the request
			return next(ctx, request)
		}
	}
}
