package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
)

var (
	EnvServiceName = os.Getenv("SERVICE_NAME")
	EnvServiceHost = os.Getenv("SERVICE_HOST")
	EnvServicePort = os.Getenv("SERVICE_PORT")
)

func registerService() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Erro ao conectar ao Consul: %v", err)
	}

	servicePort, _ := strconv.Atoi(EnvServicePort)

	registration := &api.AgentServiceRegistration{
		ID:      EnvServiceName,
		Name:    EnvServiceName,
		Address: EnvServiceHost,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", EnvServiceHost, servicePort),
			Interval: "10s",
		},
		Meta: map[string]string{},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Erro ao registrar serviço no Consul: %v", err)
	}

	log.Printf("Serviço %s registrado no Consul", EnvServiceName)
}

func main() {
	registerService()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, fmt.Sprintf("Hello from %s!", EnvServiceName))
	})

	log.Printf("%s rodando na porta %s", EnvServiceName, EnvServicePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", EnvServicePort), nil))
}
