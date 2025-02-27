package main

import (
	"encoding/json"
	"fmt"
	"github.com/yrn-go/yrn/pkg/yconnector"
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

func discoverService(serviceName string) (string, error) {
	client := getConsulClient()

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil || len(services) == 0 {
		return "", fmt.Errorf("nenhum serviço %s disponível", serviceName)
	}

	service := services[0].Service
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	serviceURL, err := discoverService(EnvConnectorServiceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	resp, err := http.Get(serviceURL)
	if err != nil {
		http.Error(w, "Erro ao chamar o serviço", http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := io.ReadAll(resp.Body)

	_, _ = w.Write(body)
}

func handlerServices(w http.ResponseWriter, r *http.Request) {
	client := getConsulClient()

	services, err := client.Agent().Services()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	schemas := map[string]any{}

	for index, service := range services {
		serviceUrl := fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, yconnector.EndpointSchema)

		resp, err := http.Get(serviceUrl)
		if err != nil {
			http.Error(w, "Erro ao chamar o serviço", http.StatusInternalServerError)
			return
		}

		resBody, _ := io.ReadAll(resp.Body)

		var respMap map[string]any
		_ = json.Unmarshal(resBody, &respMap)

		schemas[index] = respMap
	}

	_ = json.NewEncoder(w).Encode(schemas)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/services", handlerServices)

	log.Printf("Gateway rodando na porta %s", EnvServicePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", EnvServicePort), nil))
}
