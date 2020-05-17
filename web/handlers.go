package web

import (
	"encoding/json"
	"fmt"
	"github.com/google/logger"
	"github.com/kristian-d/kristiandarlington.com/docs"
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path"
)

func getTemplate(templateNames []string) (string, error) {
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
	for _, templateName := range templateNames {
		err = joinTemplate(templateName)
		if err != nil {
			return "", err
		}
	}
	return tmpl, nil
}

func (h *handler) executeTemplates(w http.ResponseWriter, templateNames []string, data interface{}) {
	rawTmpl, err := getTemplate(templateNames)
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
	err = tmpl.Execute(w, data)
}

func (h *handler) index(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.executeTemplates(w, []string{"index.html"}, data)
}

func (h *handler) projects(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.executeTemplates(w, []string{"projects.html"}, data)
}

func (h *handler) projectReports(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.logger.Info(r.URL.Path)
	filepath := r.URL.Path
	file, err := docs.Docs.Open(filepath); if err != nil {
		h.logger.Error(fmt.Sprintf("%s not found", filepath))
		http.Error(w, "failed to provide file", http.StatusNotFound)
		return
	}
	defer file.Close()
	filestats, err := file.Stat(); if err != nil {
		h.logger.Error(fmt.Sprintf("failed to get stats for %s", filepath))
		http.Error(w, "failed to provide file", http.StatusInternalServerError)
		return
	}
	_, filename := path.Split(filepath)
	modtime := filestats.ModTime()
	http.ServeContent(w, r, filename, modtime, file)
}

func (h *handler) resume(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.executeTemplates(w, []string{"resume.html"}, data)
}

func (h *handler) about(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.executeTemplates(w, []string{"about.html"}, data)
}

func (h *handler) contact(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.executeTemplates(w, []string{"contact.html", "load_animation.html"}, data)
}

func (h *handler) contactSend(w http.ResponseWriter, r *http.Request, data interface{}) {
	err := r.ParseForm(); if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnAddress := html.EscapeString(r.FormValue("email"))
	message:= html.EscapeString(r.FormValue("message"))
	recaptchaResponse := r.FormValue("g-recaptcha-response")

	if os.Getenv("ENV") == "prod" {
		validation, recaptchaErr := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
			url.Values{ "secret": {h.cfg.RecaptchaSecretKey}, "response": {recaptchaResponse}})
		if recaptchaErr != nil {
			h.logger.Error("failed to post recaptcha form to external server err=", recaptchaErr)
			http.Error(w, recaptchaErr.Error(), http.StatusInternalServerError)
			return
		}

		v := &struct{
			Success bool `json:"success"`
			ChallengeTs string `json:"challenge_ts"`
			Hostname string `json:"hostname"`
			ErrorCodes []string `json:"error-codes"`
		}{}

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

		if v.Success == false {
			h.logger.Info(fmt.Sprintf("recaptcha failed verification body=%+v\tunmarshalled=%+v", string(validationBytes), v))
			http.Error(w, "recaptcha failed verification", http.StatusInternalServerError)
			return
		}
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
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusSeeOther)
	h.contact(w, r, struct{
		CenterPopupNotification string
	}{
		CenterPopupNotification: "Message sent",
	})
}
