// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	frontend "govm/internal"
)

func handleApiIndexRequest(w http.ResponseWriter, r *http.Request) {

	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	authorization := r.Header.Get("Authorization")
	isValidSession := authorization == "Bearer "+ServerAdminSessionToken

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	var response IndexDTO
	if isValidSession {
		response = IndexDTO{
			Email:           ServerAdminEmail,
			IsAuthenticated: true,
		}
	} else {
		response = IndexDTO{
			Email:           "",
			IsAuthenticated: false,
		}
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("handleApiIndexRequest: encoding: error: %v", err)
		sendHttpError(w, EncodingFailedError, http.StatusInternalServerError)
		return
	}

}

func handleApiServerCreateRequest(w http.ResponseWriter, r *http.Request) {

	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	authorization := r.Header.Get("Authorization")
	isValidSession := authorization == "Bearer "+ServerAdminSessionToken

	if !isValidSession {
		sendHttpError(w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	// Initialize an instance of VirtualServerDTO
	var requestBody CreateVirtualServerDTO

	// Decode the request body into requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	var name string = *requestBody.Name

	if name == "" {
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	newServer := NewServerState(name, StoppedServerStatus)

	ServerCache = append(ServerCache, newServer)

	response := VirtualServerListDTO{
		Payload: convertToServerListDTO(ServerCache),
	}

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Serialize the map to JSON and write it to the response.
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("handleApiServerCreateRequest: encoding: error: %v", err)
		sendHttpError(w, EncodingFailedError, http.StatusInternalServerError)
		return
	}

}

func convertToServerDTO(
	cache ServerState,
) VirtualServerDTO {
	return VirtualServerDTO{
		Name:   cache.Name,
		Status: cache.Status,
	}
}

func convertToServerListDTO(
	cache []*ServerState,
) []VirtualServerDTO {

	// Create a slice to hold the converted DTOs
	serverDTOList := make([]VirtualServerDTO, len(cache))

	// Iterate over the slice and convert each ServerState to VirtualServerDTO
	for i, server := range cache {
		serverDTOList[i] = convertToServerDTO(*server)
	}

	return serverDTOList
}

func handleApiServerListRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		handleApiServerCreateRequest(w, r)
		return
	}

	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	authorization := r.Header.Get("Authorization")
	isValidSession := authorization == "Bearer "+ServerAdminSessionToken

	if !isValidSession {
		sendHttpError(w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	response := VirtualServerListDTO{
		Payload: convertToServerListDTO(ServerCache),
	}

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Serialize the map to JSON and write it to the response.
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("handleApiServerListRequest: encoding: error: %v", err)
		sendHttpError(w, EncodingFailedError, http.StatusInternalServerError)
		return
	}

}

func handleApiAuthRequest(w http.ResponseWriter, r *http.Request) {

	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	// Initialize an instance of VirtualServerDTO
	var requestBody AuthenticateEmailDTO

	// Decode the request body into requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	var email string = requestBody.Email
	var password string = requestBody.Password

	if email != ServerAdminEmail || password != ServerAdminPassword {
		sendHttpError(w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	token, err2 := generateAuthToken()
	if err2 != nil {
		log.Printf("handleApiAuthRequest: generating session: error: %v", err2)
		sendHttpError(w, SessionGenerationFailedError, http.StatusInternalServerError)
		return
	}

	ServerAdminSessionToken = token

	response := EmailTokenDTO{
		Token:    token,
		Email:    email,
		Verified: true,
	}

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Serialize the map to JSON and write it to the response.
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("handleApiAuthRequest: encoding: error: %v", err)
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

	r.HandleFunc("/api/v1", handleApiIndexRequest)
	r.HandleFunc("/api/v1/auth", handleApiAuthRequest)
	r.HandleFunc("/api/v1/servers", handleApiServerListRequest)
	r.Handle("/metrics", promhttp.Handler())
	r.PathPrefix("/").Handler(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(listen, r)
	if err != nil {
		panic("failed to start http server")
	}
}
