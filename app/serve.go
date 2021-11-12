package app

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed public
var contentDir embed.FS

// serveDir references the directory to be served.
// The directory is relative to this file.
const serveDir = "public"

// DevServeDir is reference to directory to be served when running in dev mode.
// This is to support hot-reloading since the dir will not be embedded in the binary.
// Directory is relative to project root.
var DevServeDir = "app/" + serveDir

// ServerRoot references the public directory since `embed` works from the current directory
// issue reference: https://stackoverflow.com/a/66248259/10813908
var ProdServeFS = func() fs.FS {
	serverRoot, err := fs.Sub(contentDir, serveDir)
	if err != nil {
		log.Fatal(err)
	}
	return serverRoot
}()
