package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippet.theolufemisamuel.com/internal/models"
	"snippet.theolufemisamuel.com/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

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
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires, err := strconv.Atoi((r.PostForm.Get("expires")))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// Initialize a map to hold any validation errors for the form fields.
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	//length of 100" and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the Valid() method to see if any of the checks failed. If they did,
	// then re-render the template passing in the form in the same way as
	// before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
