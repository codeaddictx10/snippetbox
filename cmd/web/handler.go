package main

import (
	"errors"
	"net/http"
	"strconv"

	"snippet.theolufemisamuel.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// go func() { // This will crash the server just becuase the panic is not handled in the goroutine
	// 	print("hello")
	// 	panic("oh no again")
	// }()
	// go func() { // This will not crash the server because the panic is handled in the goroutine
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			log.Print(fmt.Errorf("%s\n%s", err, debug.Stack()))
	// 		}
	// 	}()
	// 	print("hello")
	// 	panic("oh no again")
	// }()
	// panic("oh no")

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// err = ts.ExecuteTemplate(w, "base", data)

	// if err != nil {
	// 	app.serverError(w, err)
	// }
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }
	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippet: snippet,
	// }

	// err = ts.ExecuteTemplate(w, "base", data)

	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// // w.Write([]byte("Display a specific snippet"))
	// fmt.Fprintf(w, "Display a specific snippet with %d...", id)

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
