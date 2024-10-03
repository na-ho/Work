package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	goldmarkHTML "github.com/yuin/goldmark/renderer/html"
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

// Initialize Goldmark with extensions
var markdown = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("monokai"),
			highlighting.WithFormatOptions(
				html.WithLineNumbers(false),
			),
		),
	),
	goldmark.WithRendererOptions(
		goldmarkHTML.WithHardWraps(),
		goldmarkHTML.WithXHTML(),
	),
)

// renderMarkdownFile renders the Markdown file using Goldmark and serves it as HTML
func renderMarkdownFile(c *gin.Context, filePath string) {
	// Read the content of the Markdown file
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %v", err)
		return
	}

	var buf bytes.Buffer
	// Convert Markdown to HTML using the configured parser
	if err := markdown.Convert(content, &buf); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render markdown: %v", err)
		return
	}

	// Create a simple HTML page with the rendered content
	htmlContent := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>%s</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 20px; }
            /* Include CSS for syntax highlighting */
            %s
        </style>
    </head>
    <body>
        %s
    </body>
    </html>
    `, filepath.Base(filePath), highlightingCSS("monokai"), buf.String())

	// Serve the rendered Markdown as HTML
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(htmlContent))
}

// highlightingCSS returns the CSS styles for syntax highlighting
func highlightingCSS(styleName string) string {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))
	var buf bytes.Buffer
	err := formatter.WriteCSS(&buf, style)
	if err != nil {
		return ""
	}
	return buf.String()
}

func main() {
	router := gin.Default()

	// Serve the current directory
	router.Use(StaticFileMiddleware("."))

	// Start the server
	router.Run(":9000")
}
