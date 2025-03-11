package ybase

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/qri-io/jsonschema"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	MapKeySchema = "schema"

	EnvServiceName = "SERVICE_NAME"
	EnvServiceHost = "SERVICE_HOST"
	EnvServicePort = "SERVICE_PORT"

	EndpointHealth   = "/health"
	EndpointSchema   = "/" + MapKeySchema
	EndpointServices = "/services"
)

var (
	serviceName = os.Getenv(EnvServiceName)
	serviceHost = os.Getenv(EnvServiceHost)
	servicePort = os.Getenv(EnvServicePort)
)

type ServerRunFunc func() (err error)

func NewApp(
	schema *jsonschema.Schema,
	meta map[string]string,
) ServerRunFunc {
	switch "" {
	case
		serviceName,
		serviceHost,
		servicePort,
		os.Getenv(api.HTTPAddrEnvName):
		log.Panicf("envs SERVICE_NAME, SERVICE_HOST, SERVICE_PORT and %s is required\n", api.HTTPAddrEnvName)
	}

	if schema == nil {
		log.Panicln("schema is required")
	}

	registerService(meta)

	engine := gin.Default()

	engine.GET(EndpointHealth, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	engine.GET(EndpointSchema, func(c *gin.Context) {
		c.JSON(http.StatusOK, schema)
	})

	return func() (err error) {
		return engine.Run(":" + servicePort)
	}
}

func registerService(meta map[string]string) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Panicf("error connecting to Consul: %v\n", err)
	}

	servicePortInt, _ := strconv.Atoi(servicePort)

	registration := &api.AgentServiceRegistration{
		ID:      serviceName,
		Name:    serviceName,
		Address: serviceHost,
		Port:    servicePortInt,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%s/health", serviceHost, servicePort),
			Interval: "10s",
		},
		Meta: meta,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Panicf("error registering service in Consul: %v\n", err)
	}

	slog.Info(
		"service registered in Consul",
		slog.String("serviceName", serviceName),
	)
}
