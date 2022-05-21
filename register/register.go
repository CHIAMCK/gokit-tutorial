package main

import (
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
)

func Register(consulHost, consulPort, svcHost, svcPort string, logger log.Logger) (registar sd.Registrar) {

	var client consul.Client
	{
		// returns a default configuration for the client
		consulCfg := api.DefaultConfig()
		// change the default address
		consulCfg.Address = consulHost + ":" + consulPort
		// create a new consul client
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}

		client = consul.NewClient(consulClient)
	}

	// define a node or service level check
	// do health check
	check := api.AgentServiceCheck{
		// health check URL
		HTTP:     "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	// register a new service
	reg := api.AgentServiceRegistration{
		ID:      "arithmetic" + uuid.New(),
		Name:    "arithmetic",
		Address: svcHost,
		Port:    port,
		Tags:    []string{"arithmetic", "chiam"},
		Check:   &check,
	}

	// registers service instance liveness information to Consul
	registar = consul.NewRegistrar(client, &reg, logger)
	return
}

// create consul client
// create health check agent service
// register service
//
