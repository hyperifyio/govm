// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package frontend

import (
	"embed"
	"io/fs"
)

//go:embed frontend-govm/build
var FrontendFS embed.FS

var BuildFS fs.FS

func init() {

	var err error

	// Create a subdirectory in the filesystem.
	// This assumes your files are located at 'frontend-govm/build' in the embedded filesystem.
	BuildFS, err = fs.Sub(FrontendFS, "frontend-govm/build")
	if err != nil {
		panic(err) // Or handle the error as appropriate
	}

}
