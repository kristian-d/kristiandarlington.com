package web

import (
	"encoding/json"
	"github.com/kristian-d/kristiandarlington.com/config"
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
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

func contactSend(w http.ResponseWriter, r *http.Request, cfg config.Config) {
	err := r.ParseForm(); if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnAddress := r.FormValue("email")
	message:= r.FormValue("message")
	recaptchaResponse := r.FormValue("g-recaptcha-response")

	validation, recaptchaErr := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{ "secret": {cfg.RecaptchaSecretKey}, "response": {recaptchaResponse}})
	if recaptchaErr != nil {
		http.Error(w, recaptchaErr.Error(), http.StatusInternalServerError)
		return
	}

	v := struct{
		success bool `json:"success"`
		challengeTs string `json:"challenge_ts"`
		hostname string `json:"hostname"`
		errorCodes []string `json:"error-codes"`
	}{
		success: false,
		challengeTs: "",
		hostname: "",
		errorCodes: []string{},
	}

	validationBytes, validationErr := ioutil.ReadAll(validation.Body)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(validationBytes, &v); if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if v.success == false {
		http.Error(w, "recaptcha failed verification", http.StatusInternalServerError)
		return
	}

	content := "From: " + cfg.EmailCredentials.Address + "\n" +
		"To: " + cfg.EmailCredentials.Address + "\n" +
		"Subject: [Message from kristiandarlington.com]\n\n" +
		"Provided return address: " + returnAddress + "\n\n" + message

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", cfg.EmailCredentials.Address, cfg.EmailCredentials.Password, "smtp.gmail.com"),
		cfg.EmailCredentials.Address,
		[]string{cfg.EmailCredentials.Address},
		[]byte(content))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
