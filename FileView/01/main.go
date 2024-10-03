package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Custom middleware to handle directory listing and file serving with correct paths
func StaticFileMiddleware(relativePath string, fs http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedPath := c.Request.URL.Path
		println(requestedPath)
		// Ensure the path starts with the relativePath (e.g., "/logs")
		if strings.HasPrefix(requestedPath, relativePath) {
			// Remove the prefix from the requested path to get the actual file path
			//filePath := filepath.Join(".", strings.TrimPrefix(requestedPath, relativePath))
			filePath := filepath.Join(".", requestedPath)
			println(filePath)
			// Check if the path is a directory or a file
			fileInfo, err := ioutil.ReadDir(filePath)
			if err == nil {
				// If it's a directory, list the files and directories
				var fileList []string
				for _, file := range fileInfo {
					if file.IsDir() {
						// If it's a directory, ensure the link ends with a "/"
						fileList = append(fileList, fmt.Sprintf("<a href='%s/%s/'>%s/</a><br>", requestedPath, file.Name(), file.Name()))
					} else {
						// If it's a file, add a link to the file
						fileList = append(fileList, fmt.Sprintf("<a href='%s/%s'>%s</a><br>", requestedPath, file.Name(), file.Name()))
					}
				}

				// Render the list of files and directories
				c.Writer.Header().Set("Content-Type", "text/html")
				c.Writer.WriteHeader(http.StatusOK)
				templateStr := "<html><body>%s</body></html>"
				renderedTemplate := fmt.Sprintf(templateStr, strings.Join(fileList, ""))
				c.Writer.Write([]byte(renderedTemplate))
				return
			}

			// If it's not a directory, serve the file directly
			http.ServeFile(c.Writer, c.Request, filePath)
			return
		}

		c.Next() // Continue if the request path is not within /logs
	}
}

func main() {
	router := gin.Default()

	// Use the custom middleware to serve logs and directories
	router.Use(StaticFileMiddleware("/logs", http.Dir("./logs")))

	// Start the server
	router.Run(":9000")
}
