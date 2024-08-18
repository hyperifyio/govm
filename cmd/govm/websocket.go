// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		logRequest("CheckOrigin", r)
		// Allow connections from any origin
		return true
	},
}

func (api *ApiServer) onVncOpen(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		sendJsonError("onVncOpen", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}

	session := api.authenticateSession(r)
	if session == nil {
		sendJsonError("onVncOpen", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	vncPassword, err := generatePassword(8)
	if err != nil {
		logAndSendJsonError(err, "onVncOpen", w, VncGeneratePasswordError, http.StatusInternalServerError)
		return
	}

	err = api.service.SetVNCPassword(name, vncPassword)
	if err != nil {
		logAndSendJsonError(err, "onVncOpen", w, VncSetPasswordError, http.StatusInternalServerError)
		return
	}

	token, err := generatePassword(32)
	if err != nil {
		logAndSendJsonError(err, "onVncOpen", w, VncGenerateTokenError, http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("api/vnc/%s", token)
	url := fmt.Sprintf("/api/novnc/vnc_lite.html?path=%s&password=%s&scale=true", path, vncPassword)

	api.vncSessions[token] = name

	response := ServerVncDTO{
		URL:      url,
		WS:       "/" + path,
		Password: vncPassword,
		Token:    token,
	}
	sendJsonData("onVncOpen", w, response)

}

func (api *ApiServer) onVncClose(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["token"]

	if r.Method != "DELETE" {
		sendJsonError("onVncClose", w, InvalidMethod, http.StatusMethodNotAllowed)
		return
	}

	session := api.authenticateSession(r)
	if session == nil {
		log.Printf("onVncClose: Not valid session")
		sendJsonError("onVncClose", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	name, exists := api.vncSessions[token]
	if exists {
		delete(api.vncSessions, token)
	}

	vncPassword, err := generatePassword(8)
	if err != nil {
		log.Printf("onVncClose: Could not generate a VNC password: %v", err)
		return
	}

	err = api.service.SetVNCPassword(name, vncPassword)
	if err != nil {
		log.Printf("onVncClose: Could not change VNC password: %v", err)
		return
	}

	response := ServerVncDTO{}
	sendJsonData("onVncClose", w, response)

}

func (api *ApiServer) onVncWebSocket(w http.ResponseWriter, r *http.Request) {

	logRequest("onVncWebSocket", r)

	vars := mux.Vars(r)
	token := vars["token"]

	name, exists := api.vncSessions[token]
	if !exists {
		sendJsonError("onVncWebSocket", w, UnauthorizedError, http.StatusUnauthorized)
		return
	}

	vncTarget, err := api.service.GetVNC(name)
	if err != nil {
		log.Printf("onVncWebSocket: Could not get VNC: %v", err)
		return
	}
	log.Printf("onVncWebSocket: Connecting to: %s", vncTarget)

	// Upgrade the HTTP server connection to a WebSocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("onVncWebSocket: Upgrade error: %v", err)
		return
	}
	defer wsConn.Close()

	// Connect to the VNC server
	vncConn, err := net.Dial("tcp", vncTarget)
	if err != nil {
		log.Println("onVncWebSocket: Error connecting to VNC server:", err)
		return
	}
	defer vncConn.Close()

	// Forward messages from the WebSocket to the VNC server
	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("onVncWebSocket: Read error:", err)
				break
			}
			_, err = vncConn.Write(message)
			if err != nil {
				log.Println("onVncWebSocket: Write error:", err)
				break
			}
		}
	}()

	// Forward messages from the VNC server to the WebSocket
	buffer := make([]byte, 1024)
	for {
		n, err := vncConn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("onVncWebSocket: Read error:", err)
			}
			break
		}
		err = wsConn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			log.Println("onVncWebSocket: Write error:", err)
			break
		}
	}

}
