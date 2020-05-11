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
	var cfg *config.Config
	env := os.Getenv("ENV")
	switch env {
	case "prod":
		var err error
		cfg, err = config.LoadEnv()
		if err != nil {
			log.Fatal(err)
		}
	case "local":
		configPath := path.Join(projectpath.Root, "config/local.yml")
		var err error
		cfg, err = config.LoadFile(configPath)
		if err != nil { // TODO: do better
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown environment")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Index)
	mux.HandleFunc("/projects/", web.Projects)
	mux.HandleFunc("/resume/", web.Resume)
	mux.HandleFunc("/about/", web.About)
	mux.HandleFunc("/contact/", web.Contact)
	mux.Handle("/static/", http.FileServer(ui.Assets))

	mux.Handle("/.well-known/", http.FileServer(ui.Assets))

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Millisecond,
	}

	fmt.Printf("Server listening on port %s\n", cfg.Server.Port)
	log.Fatal(srv.ListenAndServe())
}
