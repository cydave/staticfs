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
	sfs := staticfs.New(static)
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

	<script src="/static/js/scripts.js"></script>
</body>
</html>
`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, page)
	})
	r.Run("127.0.0.1:3000")
}
