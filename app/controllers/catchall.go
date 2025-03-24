package controllers

import (
	"echo-fullstack/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
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
	// Start with a new template
	tmpl := template.New("page")

	// 1. Load all components (partials)
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

	// 2. Load layout and view contents
	layoutContent, err := os.ReadFile(layoutPath)
	if err != nil {
		return err
	}
	viewContent, err := os.ReadFile(viewPath)
	if err != nil {
		return err
	}

	// 3. Parse combined layout + view
	_, err = tmpl.Parse(string(layoutContent) + "\n" + string(viewContent))
	if err != nil {
		return err
	}

	// 4. Execute final template
	return tmpl.Execute(w, data)
}
