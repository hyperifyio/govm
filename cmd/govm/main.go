// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	govm "govm"
)

var ServerKey []byte

func main() {

	// Define flags
	addr := flag.String("addr", "", "change default address to listen")
	privateKeyString := flag.String("private-key", parseStringEnv("PRIVATE_KEY", ""), "set private key")
	port := flag.Int("port", parseIntEnv("PORT", 3001), "change default port")
	version := flag.Bool("version", false, "Show version information")
	initPrivateKey := flag.Bool("init-private-key", false, "Create a new private key and print it")

	listenTo := fmt.Sprintf("%s:%d", *addr, *port)

	// Parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s by %s\nURL = %s\n", govm.Name, govm.Version, govm.Author, govm.URL)
		return
	}

	if *initPrivateKey {
		key, err := generateKey()
		if err != nil {
			fmt.Printf("ERROR: Failed to generate key: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("PRIVATE_KEY=%s\n", hex.EncodeToString(key))
		}
		return
	}

	if *privateKeyString == "" {
		key, err := generateKey()
		if err != nil {
			log.Printf("ERROR: Failed to generate key: %v\n", err)
			os.Exit(1)
		} else {
			log.Printf("Warning! Initialized with a random private key '%s'. You might want to make this persistent.\n", hex.EncodeToString(key))
			ServerKey = key
		}
	} else {
		key, err := hex.DecodeString(*privateKeyString)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode private key: %v\n", err)
			os.Exit(1)
		} else {
			ServerKey = key
		}
	}

	log.Printf("Starting server at %s\n", listenTo)

	startLocalServer(listenTo)

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
