package main

import (
	"fmt"
	"github.com/zerospam/check-smtp/http/handlers"
	"github.com/zerospam/check-smtp/lib/environment-vars"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/check", handlers.CheckTransport)
	http.HandleFunc("/healthz", handlers.HealthCheck)
}

func main() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
	err := http.ListenAndServe(fmt.Sprintf(":%s", environmentvars.GetVars().ApplicationPort), nil)
	if err != nil {
		panic(err)
	}
}
