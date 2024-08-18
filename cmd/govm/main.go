// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"govm"
)

var ServerAdminEmail string
var ServerAdminPassword string

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
	features := flag.String("features", "start,stop,restart", "Enable server actions. Available actions are none, all, create, deploy, start, stop, restart, delete, and console.")

	listenTo := fmt.Sprintf("%s:%d", *addr, *port)

	// Parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s by %s\nURL = %s\n", govm.Name, govm.Version, govm.Author, govm.URL)
		return
	}

	if *adminEmail == "" {
		ServerAdminEmail = DefaultAdminUserEmail
	} else {
		ServerAdminEmail = *adminEmail
	}
	fmt.Printf("ADMIN_EMAIL=%s\n", ServerAdminEmail)

	if *adminPassword == "" {
		password, err := generatePassword(32)
		if err != nil {
			fmt.Printf("ERROR: Failed to generate admin password: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("ADMIN_PASSWORD=%s\n", password)
		}
		ServerAdminPassword = password
	} else {
		ServerAdminPassword = *adminPassword
	}

	// Features
	var err any
	var enabledActions []ServerActionCode
	featuresList := strings.Split(*features, ",")
	if contains(featuresList, "none") {
	} else if contains(featuresList, "all") {
		enabledActions = AllServerActionCodes()
	} else {
		enabledActions, err = ParseServerActionCodeList(featuresList)
	}

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

	server := NewApiServer(listenTo, service, enabledActions)

	err = server.startApiServer()
	if err != nil {
		log.Fatalf("Failed to start the API server: %v", err)
	}

	err = service.Stop()
	if err != nil {
		log.Fatalf("Failed to stop the service: %v", err)
	}

}

func parseIntEnv(key string, defaultValue int) int {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return result
}

func parseStringEnv(key string, defaultValue string) string {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	return str
}

// contains checks if a value exists in a slice of values
func contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
