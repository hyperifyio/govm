// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/rainycape/unidecode"
)

func myMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func myMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func myMin64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// containsZero checks if the slice contains a 0.
func containsZero(s []int) bool {
	for _, value := range s {
		if value == 0 {
			return true
		}
	}
	return false
}

func countZero(s []int) int {
	count := 0
	for _, value := range s {
		if value == 0 {
			count++
		}
	}
	return count
}

func MillisToISO(timestampMs int64) string {
	t := time.Unix(0, timestampMs*int64(time.Millisecond)).UTC()
	return t.Format(time.RFC3339)
}

func sanitizeName(name string) string {
	cleanName := unidecode.Unidecode(name)
	reg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	cleanName = reg.ReplaceAllString(cleanName, "")
	cleanName = strings.Join(strings.Fields(cleanName), "")
	if len(cleanName) > 32 {
		return cleanName[:32]
	}
	return strings.TrimSpace(cleanName)
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
