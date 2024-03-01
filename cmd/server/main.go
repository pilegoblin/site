package main

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	BASE_TEMPLATE = "templates/base.html"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// index route, greeting page
	r.Get("/", HandleIndex)

	r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, r, "templates/pages/projects.html", nil)
	})

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, r, "templates/pages/blog.html", nil)
	})

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("./dist/"))
		w.Header().Add("Cache-Control", "no-cache")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	fs := http.FileServer(http.Dir("./public/"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Println("Server starting...")
	http.ListenAndServe(":"+port, r)
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	//pick a random greeting
	greetings := []string{"Howdy", "Hey", "Hi"}
	greeting := greetings[rand.Intn(len(greetings))]

	RenderTemplate(w, r, "templates/pages/index.html", struct{ Greeting string }{Greeting: greeting})
}

// helper function to render the template for any page
func RenderTemplate(w http.ResponseWriter, r *http.Request, filename string, data any) {
	template := template.Must(template.ParseFiles(filename, BASE_TEMPLATE))
	var buf bytes.Buffer
	if err := template.ExecuteTemplate(&buf, "base", data); err != nil {
		return
	}

	render.HTML(w, r, buf.String())
}
