package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/gin-gonic/gin"
)

// Custom middleware to handle directory listing, syntax highlighting, and file serving
func StaticFileMiddleware(relativePath string, fs http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedPath := c.Request.URL.Path

		// Remove the prefix from the requested path to get the actual file path
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

			// Render the list of files and directories
			c.Writer.Header().Set("Content-Type", "text/html")
			c.Writer.WriteHeader(http.StatusOK)
			templateStr := "<html><body>%s</body></html>"
			renderedTemplate := fmt.Sprintf(templateStr, strings.Join(fileList, ""))
			c.Writer.Write([]byte(renderedTemplate))
			return
		}

		// Check if the file is a source code file (e.g., .go, .py, .js, etc.)
		if isCodeFile(filePath) {
			highlightCodeFile(c, filePath)
			return
		}

		// If it's a regular file, serve it directly
		http.ServeFile(c.Writer, c.Request, filePath)
	}
}

// highlightCodeFile uses Chroma to highlight source code files and render them as HTML
func highlightCodeFile(c *gin.Context, filePath string) {
	// Read the content of the source code file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %v", err)
		return
	}

	// Use Chroma to highlight the file content and render as HTML
	var highlightedContent strings.Builder
	err = quick.Highlight(&highlightedContent, string(content), filepath.Ext(filePath)[1:], "html", "monokai")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to highlight file: %v", err)
		return
	}

	// Serve the highlighted content as HTML
	tmpl := template.Must(template.New("highlighted").Parse(`
		<html>
		<head><title>Code Highlighting</title></head>
		<body><pre>{{.Content}}</pre></body>
		</html>`))

	c.Writer.Header().Set("Content-Type", "text/html")
	tmpl.Execute(c.Writer, map[string]interface{}{
		"Content": template.HTML(highlightedContent.String()), // Safely inject HTML
	})
}

// isCodeFile checks if the file is a source code file based on its extension
func isCodeFile(filePath string) bool {
	// List of common source code file extensions
	extensions := []string{".go", ".py", ".js", ".html", ".css", ".java", ".cpp", ".c", ".rb", ".php", ".sh"}

	for _, ext := range extensions {
		if strings.HasSuffix(filePath, ext) {
			return true
		}
	}
	return false
}

func main() {
	router := gin.Default()

	// Use the custom middleware to serve logs and directories
	router.Use(StaticFileMiddleware("/logs", http.Dir("./logs")))

	// Start the server
	router.Run(":9000")
}
