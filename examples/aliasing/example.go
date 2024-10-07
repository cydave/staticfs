package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cydave/staticfs"
)

//go:embed static/*
var static embed.FS

// Get root assets in the static FS. Every file that is in the top-level
// directory is returned.
func getRootAssets() []string {
	assets := make([]string, 0)
	entries, err := fs.ReadDir(static, "static")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		assets = append(assets, "/"+entry.Name())
	}
	return assets
}

func configureStaticFS(r *gin.Engine) error {
	sfs := staticfs.New(static)
	handler := sfs.Serve("/static")

	// Handle top-level static assets.
	// We alias them to point from / to /static
	// e.g.  /robots.txt => /static/robots.txt
	alias := func(to string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Request.URL.Path = "/static" + to
			r.HandleContext(c)
		}
	}
	aliases := getRootAssets()
	for _, a := range aliases {
		r.GET(a, alias(a))
	}

	// Non top-level assets are mapped as expected.
	r.GET("/static/*filepath", handler)
	return nil
}

func main() {
	r := gin.Default()
	err := configureStaticFS(r)
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		page := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>staticfs</title>
	<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
	<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
	<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
	<link rel="manifest" href="/site.webmanifest">
	<link rel="stylesheet" type="text/css" href="/static/css/styles.css">
</head>
<body>
	<h1>staticfs</h1>
	<h2>Aliased Files</h2>
	<ul>
		<li><a href="/robots.txt">/robots.txt => /static/robots.txt</a></li>
		<li><a href="/favicon.ico">/favicon.ico => /static/favicon.ico</a></li>
		<li><a href="/favicon-16x16.png">/favicon-16x16.png => /static/favicon-16x16.png</a></li>
		<li><a href="/favicon-32x32.png">/favicon-32x32.png => /static/favicon-32x32.png</a></li>
		<li><a href="/android-chrome-192x192.png">/android-chrome-192x192.png => /static/android-chrome-192x192.png</a></li>
		<li><a href="/android-chrome-512x512.png">/android-chrome-512x512.png => /static/android-chrome-512x512.png</a></li>
		<li><a href="/apple-touch-icon.png">/apple-touch-icon.png => /static/apple-touch-icon.png</a></li>
		<li><a href="/site.webmanifest">/site.webmanifest => /static/site.webmanifest</a></li>
	</ul>

	<h2>Files</h2>
	<ul>
		<li><a href="/static/css/styles.css">/static/css/styles.css</a></li>
		<li><a href="/static/js/scripts.js">/static/js/scripts.js</a></li>
	</ul>

	<script src="/static/js/scripts.js"></script>
</body>
</html>
`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, page)
	})
	r.Run("127.0.0.1:3000")
}
