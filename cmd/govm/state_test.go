// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package main

import (
	"testing"
)

func TestNewVirtualServerState(t *testing.T) {
	now := NewTimeNow()
	gs := NewVirtualServerState(now)

	if gs.Created != gs.Updated {
		t.Errorf("Expected Created (%v) and Updated (%v) to be equal", gs.Created, gs.Updated)
	}

	if gs.Created < now-1 || gs.Created > now {
		t.Errorf("Expected Created timestamp to be current, got %v", gs.Created)
	}

}
