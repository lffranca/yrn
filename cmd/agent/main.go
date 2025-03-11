package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yrn-go/yrn/pkg/ybase"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/hashicorp/consul/api"
)

var (
	EnvConnectorServiceName = os.Getenv("CONNECTOR_SERVICE_NAME")
	EnvServicePort          = os.Getenv("SERVICE_PORT")
	consulClient            *api.Client
	consulClientOnce        sync.Once
)

func getConsulClient() *api.Client {
	consulClientOnce.Do(func() {
		var err error
		consulClient, err = api.NewClient(api.DefaultConfig())
		if err != nil {
			slog.Error("Error initializing consul client: ", err)
		}
	})

	return consulClient
}

//func discoverService(serviceName string) (string, error) {
//	client := getConsulClient()
//
//	services, _, err := client.Health().Service(serviceName, "", true, nil)
//	if err != nil || len(services) == 0 {
//		return "", fmt.Errorf("nenhum serviço %s disponível", serviceName)
//	}
//
//	service := services[0].Service
//	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
//}

func handlerServiceFunc(c *gin.Context) {
	client := getConsulClient()

	services, err := client.Agent().Services()
	if err != nil {
		//http.Error(w, err.Error(), http.StatusServiceUnavailable)
		_ = c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}

	for index, service := range services {
		serviceUrl := fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, ybase.EndpointSchema)

		resp, err := http.Get(serviceUrl)
		if err != nil {
			//http.Error(w, "Erro ao chamar o serviço", http.StatusInternalServerError)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		resBody, _ := io.ReadAll(resp.Body)

		if services[index].Meta == nil {
			services[index].Meta = make(map[string]string)
		}

		services[index].Meta[ybase.MapKeySchema] = string(resBody)
	}

	c.JSON(http.StatusOK, services)
}

func main() {
	engine := gin.Default()

	engine.GET(ybase.EndpointServices, handlerServiceFunc)

	log.Fatal(engine.Run(":"+EnvServicePort), nil)
}
