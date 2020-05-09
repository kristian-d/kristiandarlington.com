package main

import (
	"fmt"
	"github.com/kristian-d/kristiandarlington.com/config"
	"github.com/kristian-d/kristiandarlington.com/internal/projectpath"
	"github.com/kristian-d/kristiandarlington.com/web"
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	configPath := path.Join(projectpath.Root, "config/kristiandarlington.com.yml")
	cfg, err := config.LoadFile(configPath)
	if err != nil { // TODO: do better
		log.Fatal("could not process config at " + configPath)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Index)
	mux.HandleFunc("/projects/", web.Projects)
	mux.HandleFunc("/resume/", web.Resume)
	mux.HandleFunc("/about/", web.About)
	mux.HandleFunc("/contact/", web.Contact)
	mux.Handle("/static/", http.FileServer(ui.Assets))

	ip := os.Getenv("IP")
	if ip == "" {
		ip = cfg.Address
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	srv := &http.Server{
		Addr:         ip + ":" + port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
	}

	fmt.Printf("Server listening on port %s\n", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
