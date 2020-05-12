package main

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/kristian-d/kristiandarlington.com/config"
	"github.com/kristian-d/kristiandarlington.com/internal/projectpath"
	"github.com/kristian-d/kristiandarlington.com/web"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	var cfg *config.Config
	lgger := logger.Init("MainLogger", true, false, ioutil.Discard)
	env := os.Getenv("ENV")
	switch env {
	case "prod":
		var err error
		cfg, err = config.LoadEnv(); if err != nil {
			lgger.Fatal(err)
		}
	case "local":
		configPath := path.Join(projectpath.Root, "config/local.yml")
		var err error
		cfg, err = config.LoadFile(configPath); if err != nil {
			lgger.Fatal(err)
		}
	default:
		lgger.Fatal(errors.New(fmt.Sprintf("unknown environment env=%s", env)))
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      web.NewRouter(cfg, lgger),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Millisecond,
	}

	lgger.Infof("Server listening on port %s\n", cfg.Server.Port)
	lgger.Fatal(srv.ListenAndServe())
}
