// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Declare a global counter
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		},
		[]string{"path"}, // Labels
	)

	failedOperationsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "govm_failed_operations_total",
			Help: "Total number of failed operations",
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		failedOperationsCounter,
	)
}

func recordFailedOperationMetric(operationName string) {
	// Increment the counter for the specific operation that failed
	failedOperationsCounter.WithLabelValues(operationName).Inc()
}
