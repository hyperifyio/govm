// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package frontend

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed frontend-govm/build
var webContent embed.FS
var BuildFrontend fs.FS

//go:embed frontend-novnc/vnc.html frontend-novnc/vnc_lite.html frontend-novnc/package.json frontend-novnc/app/* frontend-novnc/core/* frontend-novnc/vendor/*
var novncWebContent embed.FS
var BuildNoVNC fs.FS

func init() {

	var err error

	// Our frontend
	BuildFrontend, err = fs.Sub(webContent, "frontend-govm/build")
	if err != nil {
		log.Fatalf("Frontend initialization failed: %v", err)
	}

	// NoVNC source codes
	BuildNoVNC, err = fs.Sub(novncWebContent, "frontend-novnc")
	if err != nil {
		log.Fatalf("NoVNC initialization failed: %v", err)
	}

}
