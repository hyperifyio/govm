// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"os"
	"strconv"
	"strings"
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

func parseBooleanEnv(key string, defaultValue bool) bool {
	var acceptedValues = []string{
		"true",
		"on",
		"enabled",
		"1",
	}
	var str string
	if defaultValue {
		str = parseStringEnv(key, "true")
	} else {
		str = parseStringEnv(key, "false")
	}
	if contains(acceptedValues, strings.Trim(strings.ToLower(str), " ")) {
		return true
	} else {
		return false
	}
}
