package main

import (
	"log"

	"{{cookiecutter.module_path}}/internal/app"
	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/logger"
)

func main() {
	log.Println("Starting server...")

	cfg, err := config.Init("configs")
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"Name: %s, Version: %s, LogLevel: %s, LogModDev: %t, Profile: %s",
		cfg.Application.Name,
		cfg.Application.Version,
		cfg.Logger.Level,
		cfg.Logger.DevMode,
		cfg.Application.Profile,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Application.Version)
	err = app.NewApp(appLogger, cfg).Run()
	if err != nil {
		appLogger.Fatal(err)
	}

}
