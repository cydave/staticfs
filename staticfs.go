package staticfs

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type StaticFS struct {
	embRoot embed.FS
	root    http.FileSystem
	aliases []string
}

func New(static embed.FS) StaticFS {
	httpFS := http.FS(static)
	return StaticFS{embRoot: static, root: httpFS, aliases: []string{}}
}

func (s *StaticFS) Open(name string) (http.File, error) {
	f, err := s.root.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func (s *StaticFS) Serve(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file := c.Param("filepath")
		fp := filepath.Join(prefix, filepath.Clean(file))
		f, err := s.Open(fp)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		f.Close()
		http.FileServer(s.root).ServeHTTP(c.Writer, c.Request)
	}
}
