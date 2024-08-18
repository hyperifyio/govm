// Copyright (c) 2024. Hangover Games <info@hangover.games>. All rights reserved.

package main

import (
	"os"
	"strconv"
)

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
