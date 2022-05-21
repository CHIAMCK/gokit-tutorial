package main

import "errors"

// service is the container, consists of business logic
type Service interface {
	Add(a, b int) int
	Subtract(a, b int) int
	Multiply(a, b int) int
	Divide(a, b int) (int, error)

	// HealthCheck
	HealthCheck() bool
}

// implement Service interface
type ArithmeticService struct {
}

// methods
func (s ArithmeticService) Add(a, b int) int {
	return a + b
}

// Subtract implement Subtract method
func (s ArithmeticService) Subtract(a, b int) int {
	return a - b
}

// Multiply implement Multiply method
func (s ArithmeticService) Multiply(a, b int) int {
	return a * b
}

// Divide implement Divide method
func (s ArithmeticService) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("the dividend can not be zero!")
	}

	return a / b, nil
}

func (s ArithmeticService) HealthCheck() bool {
	return true
}

// function is also a type
// custom function type
type ServiceMiddleware func(Service) Service
