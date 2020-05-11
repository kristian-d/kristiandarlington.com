package web

import (
	"fmt"
	"github.com/kristian-d/kristiandarlington.com/config"
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
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

func index(w http.ResponseWriter, r *http.Request) {
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

func projects(w http.ResponseWriter, r *http.Request) {
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

func resume(w http.ResponseWriter, r *http.Request) {
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

func about(w http.ResponseWriter, r *http.Request) {
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

func contact(w http.ResponseWriter, r *http.Request) {
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

func contactSend(w http.ResponseWriter, r *http.Request, emailCreds config.EmailCredentials) {
	err := r.ParseForm(); if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	returnAddress := r.FormValue("email")
	message:= r.FormValue("message")

	content := "From: " + emailCreds.Address + "\n" +
		"To: " + emailCreds.Address + "\n" +
		"Subject: [Message from kristiandarlington.com]\n\n" +
		"Provided return address: " + returnAddress + "\n\n" + message

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", emailCreds.Address, emailCreds.Password, "smtp.gmail.com"),
		emailCreds.Address,
		[]string{emailCreds.Address},
		[]byte(content))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
