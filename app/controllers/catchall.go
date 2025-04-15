package controllers

import (
	"bytes"
	"echo-fullstack/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

func CatchAll(c echo.Context) error {
	path := c.Request().URL.Path
	if path == "/" || path == "" {
		path = "index.html"
	} else {
		path = strings.Trim(path, "/") + "/index.html"
	}

	viewPath := filepath.Join(config.RootPath(), "views", filepath.FromSlash(path))
	if _, err := os.Stat(viewPath); os.IsNotExist(err) {
		viewPath = filepath.Join(config.RootPath(), "views", "404.html")
	}

	layoutPath := filepath.Join(filepath.Dir(viewPath), "_layout.html")

	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		layoutPath = filepath.Join(config.RootPath(), "views", "_layout.html")
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	err := renderWithLayout(c.Response(), layoutPath, viewPath, map[string]interface{}{
		"Title":   strings.Title(filepath.Base(path)),
		"BaseUrl": os.Getenv("BASE_URL"),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Render Error: "+err.Error())
	}

	return nil
}

func renderWithLayout(w io.Writer, layoutPath, viewPath string, data map[string]interface{}) error {
	tmpl := template.New("page")

	// Load all components
	err := filepath.Walk(filepath.Join(config.RootPath(), "views", "components"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.Parse(string(content))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Load layout and view
	layoutContent, err := os.ReadFile(layoutPath)
	if err != nil {
		return err
	}
	viewContent, err := os.ReadFile(viewPath)
	if err != nil {
		return err
	}

	// Parse the combined template
	_, err = tmpl.Parse(string(layoutContent) + "\n" + string(viewContent))
	if err != nil {
		return err
	}

	// Execute into a buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	output := buf.String()

	// Minify HTML output if in production
	if os.Getenv("APP_ENV") == "PRODUCTION" {
		m := minify.New()
		m.AddFunc("text/html", html.Minify)
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("image/svg+xml", svg.Minify)
		m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

		minified, err := m.String("text/html", output)
		if err != nil {
			return err
		}
		output = minified
	}

	// Write output (minified or raw)
	_, err = w.Write([]byte(output))
	return err
}
