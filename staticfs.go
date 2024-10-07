package staticfs

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type StaticFS struct {
	embedFS embed.FS
	httpFS  http.FileSystem
	aliases []string
}

func New(embedfs embed.FS) *StaticFS {
	httpFS := http.FS(embedfs)
	return &StaticFS{embedFS: embedfs, httpFS: httpFS, aliases: []string{}}
}

// WithRootAliases adds top-level aliases, e.g. /static/robots.txt will be available at /robots.txt.
func (s *StaticFS) WithRootAliases() *StaticFS {
	aliases := make([]string, 0)
	entries, err := fs.ReadDir(s.embedFS, "static")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		aliases = append(aliases, "/"+entry.Name())
	}
	s.aliases = aliases
	return s
}

// Configure registers endpoints on the gin.Engine to serve static assets.
func (s *StaticFS) Configure(r *gin.Engine) {
	handler := s.serve()

	// Handle root aliases in case they are present.
	if len(s.aliases) > 0 {
		alias := func(to string) gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Request.URL.Path = "/static" + to
				r.HandleContext(c)
			}
		}
		for _, a := range s.aliases {
			r.GET(a, alias(a))
			r.HEAD(a, alias(a))
		}
	}

	r.GET("/static/*filepath", handler)
	r.HEAD("/static/*filepath", handler)
}

func (s *StaticFS) open(name string) (http.File, error) {
	f, err := s.httpFS.Open(name)
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

func (s *StaticFS) serve() gin.HandlerFunc {
	return func(c *gin.Context) {
		file := c.Param("filepath")
		fp := filepath.Join("/static", filepath.Clean(file))
		f, err := s.open(fp)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		f.Close()
		http.FileServer(s.httpFS).ServeHTTP(c.Writer, c.Request)
	}
}
