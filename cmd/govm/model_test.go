// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"testing"
)

func TestNewServerModel(t *testing.T) {
	gs := NewServerModel("testname", StartedServerStatusCode)
	if gs.Name != "testname" {
		t.Errorf("Expected Name (%v) and (%v) to be equal", gs.Name, "testname")
	}
}
