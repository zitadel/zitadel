package main

import (
	"github.com/checkr/openmock"
	"github.com/checkr/openmock/swagger_gen/restapi"
	"github.com/checkr/openmock/swagger_gen/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewOpenMockAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()
	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "OpenMock"
	parser.LongDescription = "OpenMock is a Go service that can mock services in integration tests, staging environment, or anywhere.  The goal is to simplify the process of writing mocks in various channels.  Currently it supports three channels: HTTP Kafka AMQP (e.g. RabbitMQ) The admin API allows you to manipulate the mock behaviour provided by openmock, live.  The base path for the admin API is \"/api/v1\".\n"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	om := &openmock.OpenMock{}
	om.ParseEnv()
	om.GRPCServiceMap = serviceMap
	om.GRPCEnabled = true
	om.HTTPEnabled = false
	om.TemplatesDirHotReload = false
	server.ConfigureAPI(om)
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
