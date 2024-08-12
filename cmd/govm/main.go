// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"govm"
)

var ServerAdminEmail string
var ServerAdminPassword string
var ServerAdminSessionToken string
var ServerCache []*ServerState

func main() {

	// Define flags
	addr := flag.String("addr", parseStringEnv("GOVM_ADDRESS", ""), "change default address to listen")
	adminEmail := flag.String("admin-email", parseStringEnv("GOVM_ADMIN_EMAIL", ""), "change default admin email address")
	adminPassword := flag.String("admin-password", parseStringEnv("GOVM_ADMIN_PASSWORD", ""), "change default admin password")
	port := flag.Int("port", parseIntEnv("PORT", 3001), "change default port")
	version := flag.Bool("version", false, "Show version information")

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
		password, err := generatePassword()
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

	log.Printf("Starting server at %s\n", listenTo)
	startApiServer(listenTo)

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
