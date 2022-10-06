package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/contrib/static"
)

type WebUIEndpoints struct {
	staticFS fs.FS
}

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

// func NewWebUIEndpoints(staticFS fs.FS) *WebUIEndpoints {
// 	api := &WebUIEndpoints{staticFS: staticFS}
// 	return api
// }

// func (api *WebUIEndpoints) AddRoutes(root *gin.RouterGroup) {
// 	fs := embeddedFileSystem{staticFS: api.staticFS}
// 	root.Use(static.Serve("/", fs))
// }

// func (api *WebUIEndpoints) AddRoutes(root *gin.RouterGroup) {
// 	fmt.Println("ADDING...")
// 	fs := embeddedFileSystem{staticFS: api.staticFS}
// 	root.GET("/", static.Serve("/", fs))
// 	// root.Use(static.Serve("/", fs))
// }

type embeddedFileSystem struct {
	http.FileSystem
	staticFS fs.FS
	indexes  bool
}

func (e embeddedFileSystem) Exists(prefix string, path string) bool {
	fmt.Printf("-->%s...%s\n", prefix, path)
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
