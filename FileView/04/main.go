package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/buildkite/terminal"
	"github.com/charmbracelet/glamour"
	"github.com/gin-gonic/gin"
)

func StaticFileMiddleware(root string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedPath := c.Request.URL.Path

		// Sanitize the path to prevent directory traversal attacks
		cleanPath := filepath.Clean(requestedPath)
		if strings.Contains(cleanPath, "..") {
			c.String(http.StatusForbidden, "Access denied")
			return
		}

		// Construct the full file path
		filePath := filepath.Join(root, cleanPath)

		// Get the file info
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				c.String(http.StatusNotFound, "File not found")
			} else {
				c.String(http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		if fileInfo.IsDir() {
			// List directory contents
			files, err := os.ReadDir(filePath)
			if err != nil {
				c.String(http.StatusInternalServerError, "Failed to read directory")
				return
			}

			var fileList []string
			for _, file := range files {
				name := file.Name()
				displayName := name
				if file.IsDir() {
					displayName += "/"
				}
				link := path.Join(requestedPath, name)
				fileList = append(fileList, fmt.Sprintf("<a href='%s'>%s</a><br>", link, displayName))
			}

			// Render the list of files and directories as HTML
			c.Writer.Header().Set("Content-Type", "text/html")
			c.Writer.WriteHeader(http.StatusOK)
			templateStr := "<html><body>%s</body></html>"
			renderedTemplate := fmt.Sprintf(templateStr, strings.Join(fileList, ""))
			c.Writer.Write([]byte(renderedTemplate))
			return
		}

		// If it's a Markdown file
		if strings.HasSuffix(filePath, ".md") {
			renderMarkdownFile(c, filePath)
			return
		}

		// Serve the file
		c.File(filePath)
	}
}

func ansiToHTML(ansiStr string) string {
	return string(terminal.Render([]byte(ansiStr)))
}

// renderMarkdownFile renders the Markdown file using Goldmark and serves it as HTML
func renderMarkdownFile(c *gin.Context, filePath string) {
	// Read the content of the Markdown file
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %v", err)
		return
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(), // Use default styles
		glamour.WithWordWrap(0), // No word wrapping
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create renderer: %v", err)
		return
	}

	// Render the Markdown content
	renderedANSI, err := renderer.Render(string(content))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to render markdown: %v", err)
		return
	}

	// Convert ANSI to HTML
	renderedHTML := ansiToHTML(renderedANSI)

	// Serve the rendered Markdown as HTML
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(renderedHTML))
}

func main() {
	router := gin.Default()

	// Serve the current directory
	router.Use(StaticFileMiddleware("."))

	// Start the server
	router.Run(":9000")
}
