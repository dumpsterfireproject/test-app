package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/contrib/static"
)

func EmbedFolder(fsEmbed embed.FS, targetPath string, index bool) static.ServeFileSystem {
	subFS, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embeddedFileSystem{
		FileSystem: http.FS(subFS),
		indexes:    index,
	}
}

type embeddedFileSystem struct {
	http.FileSystem
	staticFS fs.FS
	indexes  bool
}

func (e embeddedFileSystem) Exists(prefix string, path string) bool {
	f, err := e.Open(path) // e.staticFS.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// check if indexing is allowed
	s, _ := f.Stat()
	if s.IsDir() && !e.indexes {
		return false
	}

	return true
}
