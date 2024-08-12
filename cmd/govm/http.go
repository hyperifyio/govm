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

func onApiIndexRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onApiIndexRequest", r)
	isValidSession := authenticateSession(r)

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

	sendJsonData("onApiIndexRequest", w, response)
}

func onApiCreateServerRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onApiCreateServerRequest", r)

	isValidSession := authenticateSession(r)
	if !isValidSession {
		sendJsonError("onApiCreateServerRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	// Initialize an instance of ServerDTO
	var requestBody CreateServerDTO

	// Decode the request body into requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		sendJsonError("onApiCreateServerRequest", w, BadBodyError, http.StatusBadRequest)
		return
	}

	var name string = *requestBody.Name
	if name == "" {
		sendJsonError("onApiCreateServerRequest", w, BadBodyError, http.StatusBadRequest)
		return
	}

	item := NewServerState(name, StoppedServerStatusCode)
	ServerCache = append(ServerCache, item)

	response := ToServerListDTO(ServerCache)
	sendJsonData("onApiCreateServerRequest", w, response)
}

func onApiServerListRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onApiServerListRequest", r)

	isValidSession := authenticateSession(r)
	if !isValidSession {
		sendJsonError("onApiServerListRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	response := ToServerListDTO(ServerCache)
	sendJsonData("onApiServerListRequest", w, response)

}

func onApiServerRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onApiServerListRequest", r)

	vars := mux.Vars(r)
	name := vars["name"]

	isValidSession := authenticateSession(r)

	if !isValidSession {
		sendJsonError("onApiServerListRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	item, isFound := FindServerStateByName(ServerCache, name)
	if !isFound {
		sendJsonError("onApiServerListRequest", w, NotFoundError, http.StatusNotFound)
	}

	response := item.ToDTO()
	sendJsonData("onApiServerListRequest", w, response)

}

func onApiAuthRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onApiAuthRequest", r)

	// Initialize an instance of ServerDTO
	var requestBody AuthenticateEmailDTO

	// Decode the request body into requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		sendJsonError("onApiAuthRequest", w, BadBodyError, http.StatusBadRequest)
		return
	}

	var email string = requestBody.Email
	var password string = requestBody.Password

	if email != ServerAdminEmail || password != ServerAdminPassword {
		sendJsonError("onApiAuthRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	token, err2 := generateAuthToken()
	if err2 != nil {
		log.Printf("onApiAuthRequest: generating session: error: %v", err2)
		sendJsonError("onApiAuthRequest", w, SessionGenerationFailedError, http.StatusInternalServerError)
		return
	}

	ServerAdminSessionToken = token

	response := EmailTokenDTO{
		Token:    token,
		Email:    email,
		Verified: true,
	}

	sendJsonData("onApiServerListRequest", w, response)

}

func startApiServer(listen string) {

	r := mux.NewRouter()

	// Wrap the file server onr to track requests using Prometheus
	fileServerHandler := http.FileServer(http.FS(frontend.BuildFS))
	wrappedFileServerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest("wrappedFileServerHandler", r)
		fileServerHandler.ServeHTTP(w, r)
	})

	r.HandleFunc("/api/v1", onApiIndexRequest).Methods("GET")
	r.HandleFunc("/api/v1/auth", onApiAuthRequest).Methods("GET", "POST")
	r.HandleFunc("/api/v1/servers", onApiServerListRequest).Methods("GET")
	r.HandleFunc("/api/v1/servers", onApiCreateServerRequest).Methods("POST")
	r.HandleFunc("/api/v1/servers/{name}", onApiServerRequest).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())
	r.PathPrefix("/").Handler(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(listen, r)
	if err != nil {
		panic("failed to start http server")
	}
}

func logRequest(method string, r *http.Request) {
	log.Printf("%s: %s %s", method, r.Method, r.URL.Path)
	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
}

func authenticateSession(r *http.Request) bool {
	authorization := r.Header.Get("Authorization")
	return authorization == "Bearer "+ServerAdminSessionToken
}

func sendJsonError(method string, w http.ResponseWriter, code string, status int) {
	recordFailedOperationMetric(code)
	w.Header().Set("Content-Type", "application/json")
	response := ErrorDTO{
		Error: code,
		Code:  status,
	}
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("%s: encoding: error: %v", method, err)
		http.Error(w, code, status)
		return
	}
}

func sendJsonData(method string, w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("%s: encoding: error: %v", method, err)
		sendJsonError(method, w, EncodingFailedError, http.StatusInternalServerError)
		return
	}
}
