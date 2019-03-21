package main

import (
	"flag"
	"fmt"
	"github.com/VirtusLab/cloud-file-server/version"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VirtusLab/cloud-file-server/config"
	mainhandler "github.com/VirtusLab/cloud-file-server/handlers"

	"github.com/gorilla/handlers"
	"gopkg.in/yaml.v2"
)

const (
	appName = "cloud-file-server"

	defaultListenPort   = 8080
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 65 * time.Second
)

func main() {
	log.Printf("Running %s version: %s-%s", appName, version.VERSION, version.GITCOMMIT)

	configFile := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("'--config' parameter is missing")
	}

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Can not read config file: %s", err)
	}

	var cfg config.Config
	if err = yaml.UnmarshalStrict(data, &cfg); err != nil {
		log.Fatalf("Can not parse config file: %s", err)
	}

	if cfg.Listen == "" {
		cfg.Listen = fmt.Sprintf(":%d", defaultListenPort)
	}

	serverName := fmt.Sprintf("%s %s-%s", appName, version.VERSION, version.GITCOMMIT)
	handler := mainhandler.New(serverName, cfg)
	if cfg.LogRequests {
		handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
	}

	server := &http.Server{
		Addr:           cfg.Listen,
		Handler:        handler,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
		IdleTimeout:    defaultIdleTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening on %s", cfg.Listen)
	log.Fatal(server.ListenAndServe())
}
