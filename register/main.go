package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
)

func main() {

	var (
		// implements command-line flag parsing
		consulHost  = flag.String("consul.host", "", "consul ip address")
		consulPort  = flag.String("consul.port", "", "consul port")
		serviceHost = flag.String("service.host", "", "service ip address")
		servicePort = flag.String("service.port", "", "service port")
	)

	flag.Parse()

	// context is used for cancellation and data sharing (metadata in gRPC), pass req scoped values in concurrent programming
	// https://dev.to/gopher/getting-started-with-go-context-l7g
	ctx := context.Background()
	// create a channel of type error
	errChan := make(chan error)

	var logger log.Logger
	{
		// returns a logger that encodes keyvals to the Writer in logfmt format
		logger = log.NewLogfmtLogger(os.Stderr)
		// returns a new contextual logger with keyvals prepended to those passed to calls to Log
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// // gRPC
	// // initialize service, with logger as dependency
	// addservice := service.NewService(logger)
	// // initialize endpoint with service as dependency
	// addendpoint := endpoints.MakeEndpoints(addservice)
	// // create service server with all the available methods
	// grpcServer := transport.NewGRPCServer(addendpoint, logger)

	// grpcListener, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	logger.Log("during", "Listen", "err", err)
	// 	os.Exit(1)
	// }

	var svc Service
	// implement Service interface
	svc = ArithmeticService{}

	// add logging middleware
	// wrap services
	svc = LoggingMiddleware(logger)(svc)

	// return endpoint
	endpoint := MakeArithmeticEndpoint(svc)

	// add ratelimit, refill every second, set capacity 3
	// can only process 3 req per second
	// wrap endpoints
	// ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	// endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)

	// using token bucket rate limiting algo
	// a limiter controls how frequently events are allowed to happen
	rateBucket := rate.NewLimiter(rate.Every(time.Second*1), 3)
	endpoint = NewTokenBucketLimitterWithBuildIn(rateBucket)(endpoint)

	healthEndpoint := MakeHealthCheckEndpoint(svc)

	endpts := ArithmeticEndpoints{
		ArithmeticEndpoint:  endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	r := MakeHttpHandler(ctx, endpts, logger)

	// create consul client
	// create health check agent service
	// register service
	registar := Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	// create sub-routine
	go func() {
		fmt.Println("Http Server start at port:9000")
		// // start base gRPC server
		// baseServer := grpc.NewServer()
		// // register gRPC service
		// pb.RegisterMathServiceServer(baseServer, grpcServer)
		// level.Info(logger).Log("msg", "Server started successfully 🚀")
		// fmt.Println("gRPC started")
		// // start to accept gRPC request with the listener
		// baseServer.Serve(grpcListener)

		// registers instance information to a service discovery system
		registar.Register()
		handler := r
		// start the server
		// send error to the channel if it fails
		errChan <- http.ListenAndServe(":9000", handler)
	}()

	// create sub-routine
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// receive error from the channel and print it out
	// fmt.Println(<-errChan)
	error := <-errChan

	registar.Deregister()
	fmt.Println(error)
}

// gRPC
// open port to accept gRPC req
// start gRPC server, listen to certain port
// register gRPC service to base gRPC server
// initialize gRPC service server (trasport), like handler in http
// initialize endpoint
// initialize service

// http
// open port to serve http req, and attach handler
// initialize handler (transport), link endpoint to a route, decode req, encode resp,
// initialize endpoint, define the structure of req, resp. we will pass service to endpoint so that we can call it
// initialize service (business logic)

//
