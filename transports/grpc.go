package transport

import (
	"context"

	"github.com/go-kit-tutorial/endpoints"
	"github.com/go-kit-tutorial/pb"
	"github.com/go-kit/kit/log"
	gt "github.com/go-kit/kit/transport/grpc"
)

// call the endpoint
type gRPCServer struct {
	// gRPC handler
	add gt.Handler
	// for forward-compatibility, that if you changed your service files and added
	// some new methods, your binary doesn't fail if you don't implement the new
	// methods in your server.
	pb.UnimplementedMathServiceServer
}

func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) pb.MathServiceServer {
	return &gRPCServer{
		// initializes a new gRPC server, implement Handler interface
		add: gt.NewServer(
			endpoints.Add,
			decodeMathRequest,
			encodeMathResponse,
		),
	}
}

func (s *gRPCServer) Add(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	// serve gRPC request
	_, resp, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	// return response, using pb
	// we receive pb, we return pb
	return resp.(*pb.MathResponse), nil
}

// decode req
// extracts request object from an HTTP request object
// a custom function type, DecodeRequestFunc
func decodeMathRequest(_ context.Context, request interface{}) (interface{}, error) {
	// type assertion, make sure interface value holds certain type
	// use the req type in generated pb
	req := request.(*pb.MathRequest)
	// return req that we define in endpoint
	// convert pb --> endpoint
	return endpoints.MathReq{NumA: req.NumA, NumB: req.NumB}, nil
}

// encode req
func encodeMathResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.MathResp)
	// return the response type in pb
	return &pb.MathResponse{Result: resp.Result}, nil
}
