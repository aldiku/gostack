package router

import (
	"echo-fullstack/app/controllers"
	"echo-fullstack/app/middlewares"
	"echo-fullstack/config"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"log"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Init Router
func Init(app *echo.Echo) {
	renderer := &TemplateRenderer{
		templates: loadTemplates("views"),
	}
	app.Renderer = renderer
	app.Use(middlewares.Cors())
	app.Use(middlewares.Gzip())
	app.Use(middlewares.Logger())
	app.Use(middlewares.Secure())
	app.Use(middlewares.Recover())

	v1 := app.Group("v1", middlewares.StripHTMLMiddleware)
	{
		v1.GET("/product", controllers.Product)
	}

	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.GET("/swagger/doc.json", func(c echo.Context) error {
		return c.File("docs/swagger.json")
	})
	app.GET("/docs", func(c echo.Context) error {
		tmplBytes, err := os.ReadFile(filepath.Join(config.RootPath(), "views", filepath.FromSlash("docs.html")))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Template Parse Error: "+err.Error())
		}

		tmpl, err := template.New("docs.html").Parse(string(tmplBytes))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Template Parse Error: "+err.Error())
		}

		return tmpl.Execute(c.Response(), map[string]interface{}{
			"Title":   "Api Documentation " + config.LoadConfig().AppName,
			"BaseUrl": os.Getenv("BASE_URL"),
		})
	})
	app.Static("/assets", config.RootPath()+"/assets")
	app.GET("/home", controllers.Home)
	app.GET("/*", controllers.CatchAll)

	log.Printf("Server started...")
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func loadTemplates(root string) *template.Template {
	tmpl := template.New("")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			name := filepath.ToSlash(relPath) // make Windows-friendly
			_, err = tmpl.New(name).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	// DEBUG
	if os.Getenv("ENVIRONMENT") != "PRODUCTION" {
		for _, t := range tmpl.Templates() {
			fmt.Println("Parsed template:", t.Name())
		}
	}

	return tmpl
}
