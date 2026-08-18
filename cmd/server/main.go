package main

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	BaseTemplate = "templates/base.html"
	IndexPage    = "templates/pages/index.html"
	ProjectsPage = "templates/pages/projects.html"
)

// available to every template, so pages that render with nil data can use them too
var templateFuncs = template.FuncMap{
	"year": func() int { return time.Now().Year() },
}

// every page is parsed once at startup instead of on each request, which used
// to re-read and re-parse both files from disk for every hit. each page gets
// its own template set because they all define their own title/content blocks.
var templates = map[string]*template.Template{
	IndexPage:    parsePage(IndexPage),
	ProjectsPage: parsePage(ProjectsPage),
}

func parsePage(filename string) *template.Template {
	return template.Must(
		template.New(filepath.Base(BaseTemplate)).Funcs(templateFuncs).ParseFiles(filename, BaseTemplate),
	)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// index route, greeting page
	r.Get("/", HandleIndex)

	r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, r, ProjectsPage, nil)
	})

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("./dist/"))
		w.Header().Add("Cache-Control", "no-cache")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	fs := http.FileServer(http.Dir("./public/"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// RFC 9116 vulnerability disclosure endpoint, served from public/.well-known/
	r.Handle("/.well-known/*", fs)

	// crawler directives, served from public/robots.txt
	r.Handle("/robots.txt", fs)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

type indexArgs struct {
	Greeting string
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	// pick a random greeting
	greetings := []string{"Howdy", "Hey", "Hi"}
	greeting := greetings[rand.Intn(len(greetings))]

	// args
	args := indexArgs{Greeting: greeting}

	RenderTemplate(w, r, IndexPage, args)
}

// helper function to render the template for any page
func RenderTemplate(w http.ResponseWriter, r *http.Request, filename string, data any) {
	tmpl, ok := templates[filename]
	if !ok {
		log.Printf("no template parsed for %s", filename)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// render into a buffer first so a half-written page never reaches the client
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		log.Printf("rendering %s: %v", filename, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	render.HTML(w, r, buf.String())
}
