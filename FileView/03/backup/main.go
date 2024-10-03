package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
)

// Custom middleware to handle directory listing, Markdown rendering, and file serving
func StaticFileMiddleware(relativePath string, fs http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedPath := c.Request.URL.Path

		// Remove the prefix from the requested path to get the actual file path (start at parent of /logs)
		filePath := filepath.Join(".", strings.TrimPrefix(requestedPath, relativePath))

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

			// Render the list of files and directories as HTML
			c.Writer.Header().Set("Content-Type", "text/html")
			c.Writer.WriteHeader(http.StatusOK)
			templateStr := "<html><body>%s</body></html>"
			renderedTemplate := fmt.Sprintf(templateStr, strings.Join(fileList, ""))
			c.Writer.Write([]byte(renderedTemplate))
			return
		}

		// If it's not a directory, check if it's a Markdown file
		if strings.HasSuffix(filePath, ".md") {
			renderMarkdownFile(c, filePath)
			return
		}

		// If it's a regular file, serve it directly
		http.ServeFile(c.Writer, c.Request, filePath)
	}
}

// renderMarkdownFile renders the Markdown file using Glamour and serves it as HTML
func renderMarkdownFile(c *gin.Context, filePath string) {
	fmt.Println("start renderMarkdownFile")
	// Read the content of the Markdown file
	/*
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read file: %v", err)
			return
		}

		// Render the Markdown content to HTML using Glamour
		rendered, err := glamour.Render(string(content), "dark")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to render markdown: %v", err)
			return
		}
	*/
	content := `# Hello World

	This is a simple example of Markdown rendering with Glamour!
	Check out the [other examples](https://github.com/charmbracelet/glamour/tree/master/examples) too.

	Bye!
	`

	var buf bytes.Buffer
	// Convert Markdown to HTML
	if err := goldmark.Convert(content, &buf); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render markdown: %v", err)
		return
	}

	// Serve the rendered Markdown as HTML
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write(buf.Bytes())
}

func main() {
	router := gin.Default()

	// Serve the parent directory of /logs
	router.Use(StaticFileMiddleware("/logs", http.Dir(".")))

	// Start the server
	router.Run(":9000")
}
