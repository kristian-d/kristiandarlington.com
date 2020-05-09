package web

import (
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

func getTemplate(templateName string) (string, error) {
	var tmpl string

	joinTemplate := func(name string) error {
		file, err := ui.Assets.Open(path.Join("/templates", name))
		if err != nil {
			return err
		}
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		tmpl += string(bytes)
		return nil
	}

	err := joinTemplate("_base.html")
	if err != nil {
		return "", err
	}
	err = joinTemplate(templateName)
	if err != nil {
		return "", err
	}
	return tmpl, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate( "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func Projects(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("projects.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func Resume(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("resume.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("contact.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}
