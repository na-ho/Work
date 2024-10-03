package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

// WebDAVMiddleware serves WebDAV requests and also supports web-based browsing.
func WebDAVMiddleware(relativePath string, fs webdav.FileSystem) gin.HandlerFunc {
	webdavHandler := &webdav.Handler{
		Prefix:     relativePath,
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(), // In-memory lock system for WebDAV operations
	}

	return func(c *gin.Context) {
		// Serve web browsing if it's a GET request and requested by a browser
		if strings.HasPrefix(c.Request.URL.Path, relativePath) {
			if c.Request.Method == http.MethodGet {
				// Serve web-based browsing for GET requests
				serveWebDirectory(c, strings.TrimPrefix(c.Request.URL.Path, relativePath), fs)
				return
			}

			// For non-GET methods (WebDAV methods like PROPFIND, PUT, DELETE), handle via WebDAV
			webdavHandler.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		// Continue to the next middleware if not a WebDAV request
		c.Next()
	}
}

// serveWebDirectory serves an HTML page that lists the files and directories in the requested path.
func serveWebDirectory(c *gin.Context, relativePath string, fs webdav.FileSystem) {
	fullPath := filepath.Join(".", relativePath) // Map relative path to full filesystem path

	// Read directory contents
	entries, err := ioutil.ReadDir(fullPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read directory: %s", err.Error())
		return
	}

	// Build HTML template to list files and directories
	var fileList []string
	for _, entry := range entries {
		if entry.IsDir() {
			fileList = append(fileList, fmt.Sprintf("<a href='%s/'>%s/</a><br>", entry.Name(), entry.Name()))
		} else {
			fileList = append(fileList, fmt.Sprintf("<a href='%s'>%s</a><br>", entry.Name(), entry.Name()))
		}
	}

	// Simple HTML page template
	tmpl := template.Must(template.New("dirlist").Parse(`
		<html>
		<head><title>Directory Listing</title></head>
		<body>
		<h1>Directory Listing for {{.Path}}</h1>
		{{.FileList}}
		</body>
		</html>
	`))

	// Serve the HTML page
	c.Writer.Header().Set("Content-Type", "text/html")
	tmpl.Execute(c.Writer, map[string]interface{}{
		"Path":     relativePath,
		"FileList": template.HTML(strings.Join(fileList, "")), // Safely join HTML content
	})
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Use the WebDAV middleware to serve files from ./files under /webdav
	router.Use(WebDAVMiddleware("/webdav", webdav.Dir("./files")))

	// Start the Gin server on port 8080
	router.Run(":8080")
}
