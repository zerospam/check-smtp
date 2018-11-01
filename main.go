package main

import (
	"fmt"
	firewall_handlers "github.com/zerospam/check-firewall/lib/handlers"
	"github.com/zerospam/check-smtp/http/handlers"
	"github.com/zerospam/check-smtp/lib/environment-vars"
	"net/http"
)

func init() {
	http.HandleFunc("/check", handlers.CheckTransport)
	http.HandleFunc("/healthz", firewall_handlers.HealthCheck)
}

func main() {

	err := http.ListenAndServe(fmt.Sprintf(":%s", environmentvars.GetVars().ApplicationPort), nil)
	if err != nil {
		panic(err)
	}
}
