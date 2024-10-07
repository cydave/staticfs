package main

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cydave/staticfs"
)

//go:embed static/*
var static embed.FS

func main() {
	r := gin.Default()
	// Set caching headers for resources that are found.
	okCallback := func(c *gin.Context, path string) {
		c.Header("Cache-Control", "private, max-age=3600")
	}
	// Set no-cache headers for resources that were not found.
	errCallback := func(c *gin.Context, err error) {
		c.Header("Pragma", "no-cache")
		c.Header("Cache-Control", "private, no-cache")
	}
	// Create staticfs with the above defined callbacks.
	sfs := staticfs.New(static).WithOKCallback(okCallback).WithErrCallback(errCallback)
	sfs.Configure(r)

	r.GET("/", func(c *gin.Context) {
		page := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>staticfs</title>
	<link rel="stylesheet" type="text/css" href="/static/css/styles.css">
</head>
<body>
	<h1>staticfs</h1>
	<h2>Files</h2>
	<ul>
		<li><a href="/static/css/styles.css">/static/css/styles.css</a></li>
		<li><a href="/static/js/scripts.js">/static/js/scripts.js</a></li>
	</ul>
	<h2>Caching</h2>
	<p>Resources that are found are cached for 3600 seconds; resources that are not found have an explicit no-cache header.</p>

	<script src="/static/js/scripts.js"></script>
</body>
</html>
`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, page)
	})
	r.Run("127.0.0.1:3000")
}
