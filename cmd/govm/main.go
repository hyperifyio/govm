// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"govm"
)

func main() {

	// Define flags
	addr := flag.String("addr", parseStringEnv("GOVM_ADDRESS", ""), "change default address to listen")
	system := flag.String("system", parseStringEnv("GOVM_SYSTEM", "qemu:///system"), "change default virtio system")
	ifType := flag.String("default-if", parseStringEnv("GOVM_INTERFACE_TYPE", "network"), "change default virtio interface type (user or network)")
	ifNetworkName := flag.String("default-if-network", parseStringEnv("GOVM_INTERFACE_NETWORK", "default"), "change default virtio network name (if network type)")
	defaultBridge := flag.String("default-bridge", parseStringEnv("GOVM_BRIDGE", "br0"), "change default virtio network bridge interface")
	volumesDir := flag.String("volumes", parseStringEnv("GOVM_VOLUMES", "./volumes"), "change default location for volumes")
	imagesDir := flag.String("images", parseStringEnv("GOVM_IMAGES", "./images"), "change default location for images")
	adminEmail := flag.String("admin-email", parseStringEnv("GOVM_ADMIN_EMAIL", ""), "change default admin email address")
	adminPassword := flag.String("admin-password", parseStringEnv("GOVM_ADMIN_PASSWORD", ""), "change default admin password")
	port := flag.Int("port", parseIntEnv("PORT", 3001), "change default port")
	version := flag.Bool("version", false, "Show version information")
	demo := flag.Bool("demo", false, "Use demo version of the service")
	features := flag.String("features", parseStringEnv("GOVM_FEATURES", "start,stop,restart,console"), "Enable server actions. Available actions are none, all, create, deploy, start, stop, restart, delete, and console.")
	configFile := flag.String("config", parseStringEnv("GOVM_CONFIG", "./config.yml"), "Configuration file")
	https := flag.Bool("https", parseBooleanEnv("GOVM_HTTPS", false), "Enable HTTPS instead of HTTP")
	certDir := flag.String("cert-dir", parseStringEnv("GOVM_CERT_DIR", "./certs"), "TLS files for HTTPS")
	certFile := flag.String("cert", parseStringEnv("GOVM_CERT_FILE", "./server.crt"), "Certificate file for HTTPS")
	keyFile := flag.String("key", parseStringEnv("GOVM_KEY_FILE", "./server.key"), "Key file for HTTPS")

	listenTo := fmt.Sprintf("%s:%d", *addr, *port)

	// Parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s by %s\nURL = %s\n", govm.Name, govm.Version, govm.Author, govm.URL)
		return
	}

	var serverAdminEmail string
	if *adminEmail == "" {
		serverAdminEmail = DefaultAdminUserEmail
	} else {
		serverAdminEmail = *adminEmail
	}
	fmt.Printf("ADMIN_EMAIL=%s\n", serverAdminEmail)

	var serverAdminPassword string
	if *adminPassword == "" {
		password, err := generatePassword(32)
		if err != nil {
			fmt.Printf("ERROR: Failed to generate admin password: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("ADMIN_PASSWORD=%s\n", password)
		}
		serverAdminPassword = password
	} else {
		serverAdminPassword = *adminPassword
	}

	var err any

	// Config
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %s: %v", *configFile, err)
	}
	configManager := NewConfigManager(*configFile, config)

	// Features
	var enabledActions []ServerActionCode
	featuresList := strings.Split(*features, ",")
	if contains(featuresList, "none") {
	} else if contains(featuresList, "all") {
		enabledActions = AllServerActionCodes()
	} else {
		enabledActions, err = ParseServerActionCodeList(featuresList)
	}

	// AuthorizationService
	authorizationService := NewSingleMemoryAuthorizationService(serverAdminEmail, serverAdminPassword)

	// SessionService
	sessionService := NewSingleMemorySessionService()

	// Service
	var service ServerService
	if *demo {
		service = NewDummyService()
		log.Printf("Starting dummy server at %s\n", listenTo)
	} else {

		absImagesDir, err := filepath.Abs(*imagesDir)
		if err != nil {
			log.Fatalf("Failed to get absolute path for images directory: %s: %v", *imagesDir, err)
		}

		absVolumesDir, err := filepath.Abs(*volumesDir)
		if err != nil {
			log.Fatalf("Failed to get absolute path for volumes directory: %s: %v", *volumesDir, err)
		}

		service = NewVirtioService(*system, absImagesDir, absVolumesDir, *ifType, *ifNetworkName, *defaultBridge, enabledActions)
		log.Printf("Starting virtio server at %s\n", listenTo)
	}

	err = service.Start()
	if err != nil {
		log.Fatalf("Failed to start the service: %v", err)
	}

	tlsEnabled := *https
	tlsDir := *certDir
	tlsCertFile := filepath.Join(tlsDir, *certFile)
	tlsKeyFile := filepath.Join(tlsDir, *keyFile)
	if tlsEnabled {
		log.Printf("Using HTTPS with certificate file (%s) and key file (%s)", tlsCertFile, tlsKeyFile)
	} else {
		log.Printf("Warning! Using unsecured HTTP")
	}

	server := NewApiServer(listenTo, tlsEnabled, tlsCertFile, tlsKeyFile, service, sessionService, authorizationService, enabledActions, configManager)

	err = server.startApiServer()
	if err != nil {
		log.Fatalf("Failed to start the API server: %v", err)
	}

	err = service.Stop()
	if err != nil {
		log.Fatalf("Failed to stop the service: %v", err)
	}

}
