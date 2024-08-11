// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"govm/internal/frontend"
)

func handleVirtualServerStateRequest(w http.ResponseWriter, r *http.Request) {

	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	// Initialize an instance of VirtualServerDTO
	var vmRequest CreateVirtualServerDTO

	// Decode the request body into vmRequest
	err := json.NewDecoder(r.Body).Decode(&vmRequest)
	if err != nil {
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	var requestedName string
	if vmRequest.Name != nil {
		requestedName = *(vmRequest.Name)
	} else {
		requestedName = ""
	}

	response := VirtualServerDTO{
		Name: requestedName,
	}

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Serialize the map to JSON and write it to the response.
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("handleVirtualServerStateRequest: encoding: error: %v", err)
		sendHttpError(w, EncodingFailedError, http.StatusInternalServerError)
		return
	}

}

func startLocalServer(listen string) {

	r := mux.NewRouter()

	// Wrap the file server handler to track requests using Prometheus
	fileServerHandler := http.FileServer(http.FS(frontend.BuildFS))
	wrappedFileServerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
		fileServerHandler.ServeHTTP(w, r)
	})

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/api/v1", handleVirtualServerStateRequest)
	r.PathPrefix("/").Handler(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(listen, r)
	if err != nil {
		panic("failed to start http server")
	}
}
