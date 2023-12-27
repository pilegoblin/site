package main

import (
	"bytes"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	BASE_TEMPLATE = "templates/base.html"
)

func main() {
	// boilerplate
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// index route, greeting page
	r.Get("/", HandleIndex)
	// temp page at the moment
	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, r, "templates/pages/about.html", nil)
	})
	// fileserver for css generated by tailwind
	staticFileServer := http.FileServer(http.Dir("./dist/"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFileServer))

	publicFileServer := http.FileServer(http.Dir("./public/"))
	r.Handle("/public/*", http.StripPrefix("/public/", publicFileServer))

	//turn that server on
	http.ListenAndServe(":3000", r)
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
		//todo: throw up a 504 page here or something
		return
	}

	render.HTML(w, r, buf.String())
}
