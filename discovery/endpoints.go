package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

func MakeDiscoverEndpoint(ctx context.Context, client consul.Client, logger log.Logger) endpoint.Endpoint {
	serviceName := "arithmetic"
	tags := []string{"arithmetic", "chiam"}
	passingOnly := true
	duration := 500 * time.Millisecond

	// 基于consul客户端、服务名称、服务标签等信息，
	// 创建consul的连接实例 (Connection instance)，可实时查询服务实例的状态信息
	// returns a Consul instancer that publishes instances for the requested service
	// need to access all the registered instances for the requested service
	instancer := consul.NewInstancer(client, logger, serviceName, tags, passingOnly)

	// 针对calculate接口创建sd.Factory
	// converts an instance string (e.g. host:port) to a specific endpoint.
	// Instances that provide multiple endpoints require multiple factories
	factory := arithmeticFactory(ctx, "POST", "calculate")

	// 使用consul连接实例（发现服务系统）、factory创建sd.Factory
	// creates an Endpointer that subscribes to updates from Instancer src and
	// uses factory f to create Endpoints.
	endpointer := sd.NewEndpointer(instancer, factory, logger)

	// 创建RoundRibbon负载均衡器
	// returns a load balancer that returns services in sequence
	// might have more than 1 instance
	balancer := lb.NewRoundRobin(endpointer)

	// 为负载均衡器增加重试功能，同时该对象为endpoint.Endpoint
	// retry if error occurs
	retry := lb.Retry(1, duration, balancer)

	return retry
}

// initialize endpoint, consul, factory, LB
// link route with endpoint
