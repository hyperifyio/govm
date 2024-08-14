// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	frontend "govm/internal"
)

type ApiServer struct {
	r                         *mux.Router
	listen                    string
	authenticatedSessionToken string
	service                   ServerService
}

func NewApiServer(listen string, service ServerService) *ApiServer {
	return &ApiServer{
		listen:                    listen,
		authenticatedSessionToken: "",
		service:                   service,
	}
}

func (api *ApiServer) onIndexRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onIndexRequest", r)
	isValidSession := api.authenticateSession(r)

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

	sendJsonData("onIndexRequest", w, response)
}

func (api *ApiServer) onAddServerRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onAddServerRequest", r)

	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onAddServerRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	// Initialize an instance of ServerDTO
	var requestBody CreateServerDTO

	// Decode the request body into requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logAndSendJsonError(err, "onAddServerRequest", w, BadBodyError, http.StatusBadRequest)
		return
	}

	var name string = *requestBody.Name
	if name == "" {
		sendJsonError("onAddServerRequest", w, BadBodyError, http.StatusBadRequest)
		return
	}

	_, err = api.service.AddServer(name)
	if err != nil {
		logAndSendJsonError(err, "onAddServerRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}

	serverList, err := api.service.GetServerList()
	if err != nil {
		logAndSendJsonError(err, "onAddServerRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}

	response := ToServerListDTO(serverList)
	sendJsonData("onAddServerRequest", w, response)
}

func (api *ApiServer) onServerListRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onServerListRequest", r)

	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onServerListRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	serverList, err := api.service.GetServerList()
	if err != nil {
		logAndSendJsonError(err, "onServerListRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}

	response := ToServerListDTO(serverList)
	sendJsonData("onServerListRequest", w, response)

}

func (api *ApiServer) onServerRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onServerListRequest", r)

	vars := mux.Vars(r)
	name := vars["name"]

	isValidSession := api.authenticateSession(r)

	if !isValidSession {
		sendJsonError("onServerListRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	item, err := api.service.FindServer(name)
	if err != nil {
		logAndSendJsonError(err, "onServerListRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if item == nil {
		sendJsonError("onServerListRequest", w, NotFoundError, http.StatusNotFound)
	} else {
		response := item.ToDTO()
		sendJsonData("onServerListRequest", w, response)
	}

}

func (api *ApiServer) onServerDeployRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onServerDeployRequest", r)
	if r.Method != "POST" {
		sendJsonError("onServerDeployRequest", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	name := vars["name"]
	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onServerDeployRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}
	item, err := api.service.DeployServer(name)
	if err != nil {
		logAndSendJsonError(err, "onServerDeployRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if item == nil {
		sendJsonError("onServerDeployRequest", w, NotFoundError, http.StatusNotFound)
	} else {
		response := item.ToDTO()
		sendJsonData("onServerDeployRequest", w, response)
	}
}

func (api *ApiServer) onServerStartRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onServerStartRequest", r)
	if r.Method != "POST" {
		sendJsonError("onServerStartRequest", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	name := vars["name"]
	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onServerStartRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}
	item, err := api.service.StartServer(name)
	if err != nil {
		logAndSendJsonError(err, "onServerStartRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if item == nil {
		sendJsonError("onServerStartRequest", w, NotFoundError, http.StatusNotFound)
	} else {
		response := item.ToDTO()
		sendJsonData("onServerStartRequest", w, response)
	}
}

func (api *ApiServer) onServerStopRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onServerStopRequest", r)
	if r.Method != "POST" {
		sendJsonError("onServerStopRequest", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	name := vars["name"]
	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onServerStopRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}
	item, err := api.service.StopServer(name)
	if err != nil {
		logAndSendJsonError(err, "onServerStopRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if item == nil {
		sendJsonError("onServerStopRequest", w, NotFoundError, http.StatusNotFound)
	} else {
		response := item.ToDTO()
		sendJsonData("onServerStopRequest", w, response)
	}
}

func (api *ApiServer) onServerRestartRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onServerRestartRequest", r)
	if r.Method != "POST" {
		sendJsonError("onServerRestartRequest", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	name := vars["name"]
	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onServerRestartRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}
	item, err := api.service.RestartServer(name)
	if err != nil {
		logAndSendJsonError(err, "onServerRestartRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if item == nil {
		sendJsonError("onServerRestartRequest", w, NotFoundError, http.StatusNotFound)
	} else {
		response := item.ToDTO()
		sendJsonData("onServerRestartRequest", w, response)
	}
}

func (api *ApiServer) onServerDeleteRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("onServerDeleteRequest", r)
	if r.Method != "POST" {
		sendJsonError("onServerDeleteRequest", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	name := vars["name"]
	isValidSession := api.authenticateSession(r)
	if !isValidSession {
		sendJsonError("onServerDeleteRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}
	item, err := api.service.DeleteServer(name)
	if err != nil {
		logAndSendJsonError(err, "onServerDeleteRequest", w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if item == nil {
		sendJsonError("onServerDeleteRequest", w, NotFoundError, http.StatusNotFound)
	} else {
		response := item.ToDTO()
		sendJsonData("onServerDeleteRequest", w, response)
	}
}

func (api *ApiServer) onAuthRequest(w http.ResponseWriter, r *http.Request) {

	logRequest("onAuthRequest", r)

	// Initialize an instance of ServerDTO
	var requestBody AuthenticateEmailDTO

	// Decode the request body into requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logAndSendJsonError(err, "onAuthRequest", w, BadBodyError, http.StatusBadRequest)
		return
	}

	var email string = requestBody.Email
	var password string = requestBody.Password

	if email != ServerAdminEmail || password != ServerAdminPassword {
		sendJsonError("onAuthRequest", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	token, err2 := generateAuthToken()
	if err2 != nil {
		log.Printf("onAuthRequest: generating session: error: %v", err2)
		sendJsonError("onAuthRequest", w, SessionGenerationFailedError, http.StatusInternalServerError)
		return
	}

	api.authenticatedSessionToken = token

	response := EmailTokenDTO{
		Token:    token,
		Email:    email,
		Verified: true,
	}

	sendJsonData("onServerListRequest", w, response)

}

func (api *ApiServer) startApiServer() error {

	api.r = mux.NewRouter()

	// Wrap the file server onr to track requests using Prometheus
	fileServerHandler := http.FileServer(http.FS(frontend.BuildFS))
	wrappedFileServerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest("wrappedFileServerHandler", r)
		fileServerHandler.ServeHTTP(w, r)
	})

	api.r.HandleFunc("/api/v1", api.onIndexRequest).Methods("GET")
	api.r.HandleFunc("/api/v1/auth", api.onAuthRequest).Methods("GET", "POST")
	api.r.HandleFunc("/api/v1/servers", api.onServerListRequest).Methods("GET")
	api.r.HandleFunc("/api/v1/servers", api.onAddServerRequest).Methods("POST")
	api.r.HandleFunc("/api/v1/servers/{name}", api.onServerRequest).Methods("GET")
	api.r.HandleFunc("/api/v1/servers/{name}/deploy", api.onServerDeployRequest).Methods("GET", "POST")
	api.r.HandleFunc("/api/v1/servers/{name}/start", api.onServerStartRequest).Methods("GET", "POST")
	api.r.HandleFunc("/api/v1/servers/{name}/stop", api.onServerStopRequest).Methods("GET", "POST")
	api.r.HandleFunc("/api/v1/servers/{name}/restart", api.onServerRestartRequest).Methods("GET", "POST")
	api.r.HandleFunc("/api/v1/servers/{name}/delete", api.onServerDeleteRequest).Methods("GET", "POST")
	api.r.Handle("/metrics", promhttp.Handler())
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(api.listen, api.r)
	if err != nil {
		return fmt.Errorf("failed to start http server: %v", err)
	}
	return nil
}

func (api *ApiServer) authenticateSession(r *http.Request) bool {
	if api.authenticatedSessionToken == "" {
		return false
	}
	authorization := r.Header.Get("Authorization")
	return authorization == "Bearer "+api.authenticatedSessionToken
}

func logRequest(method string, r *http.Request) {
	log.Printf("%s: %s %s", method, r.Method, r.URL.Path)
	httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
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
		log.Printf("%s: ERROR: encoding: %v", method, err)
		http.Error(w, code, status)
		return
	}
}

func logAndSendJsonError(err any, method string, w http.ResponseWriter, code string, status int) {
	log.Printf("%s: ERROR: %v", method, err)
	sendJsonError(method, w, code, status)
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
