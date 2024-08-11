// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package frontend

import (
	"embed"
	"io/fs"
)

//go:embed frontend-govm/frontend/build
var FrontendFS embed.FS

var BuildFS fs.FS

func init() {

	var err error

	// Create a subdirectory in the filesystem.
	// This assumes your files are located at 'project-govm/frontend/build' in the embedded filesystem.
	BuildFS, err = fs.Sub(FrontendFS, "project-govm/frontend/build")
	if err != nil {
		panic(err) // Or handle the error as appropriate
	}

}
