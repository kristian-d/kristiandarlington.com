package web

import (
	"encoding/json"
	"fmt"
	"github.com/google/logger"
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

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("index.html")
	if err != nil {
		h.logger.Error("failed to get template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		h.logger.Error("failed to parse template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func (h *handler) projects(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("projects.html")
	if err != nil {
		h.logger.Error("failed to get template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		h.logger.Error("failed to parse template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func (h *handler) resume(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("resume.html")
	if err != nil {
		h.logger.Error("failed to get template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		h.logger.Error("failed to parse template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func (h *handler) about(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("about.html")
	if err != nil {
		h.logger.Error("failed to get template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		h.logger.Error("failed to parse template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func (h *handler) contact(w http.ResponseWriter, r *http.Request) {
	rawTmpl, err := getTemplate("contact.html")
	if err != nil {
		h.logger.Error("failed to get template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("").Parse(rawTmpl)
	if err != nil {
		h.logger.Error("failed to parse template err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func (h *handler) contactSend(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm(); if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnAddress := r.FormValue("email")
	message:= r.FormValue("message")
	recaptchaResponse := r.FormValue("g-recaptcha-response")

	validation, recaptchaErr := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{ "secret": {h.cfg.RecaptchaSecretKey}, "response": {recaptchaResponse}})
	if recaptchaErr != nil {
		h.logger.Error("failed to post recaptcha form to external server err=", recaptchaErr)
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
		h.logger.Error("failed to read recaptcha response err=", validationErr)
		http.Error(w, validationErr.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(validationBytes, &v); if err != nil {
		h.logger.Error("failed to unmarshal recaptcha response err=", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if v.success == false {
		h.logger.Info(fmt.Sprintf("recaptcha failed verification body=%+v\tunmarshalled=%+v", string(validationBytes), v))
		http.Error(w, "recaptcha failed verification", http.StatusInternalServerError)
		return
	}

	content := "From: " + h.cfg.EmailCredentials.Address + "\n" +
		"To: " + h.cfg.EmailCredentials.Address + "\n" +
		"Subject: [Message from kristiandarlington.com]\n\n" +
		"Provided return address: " + returnAddress + "\n\n" + message

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", h.cfg.EmailCredentials.Address, h.cfg.EmailCredentials.Password, "smtp.gmail.com"),
		h.cfg.EmailCredentials.Address,
		[]string{h.cfg.EmailCredentials.Address},
		[]byte(content))
	if err != nil {
		logger.Error("failed to send mail to contact email err=", err, "\tcontent=", content)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
