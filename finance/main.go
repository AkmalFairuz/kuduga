package main

import (
	"github.com/akmalfairuz/finance/module/envutils"
	"github.com/akmalfairuz/finance/server"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	_ = os.Setenv("TZ", "UTC")

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	var cfg server.Config
	if err := envutils.Load(".env", &cfg); err != nil {
		log.Fatalf("error loading config: %+v", err)
	}

	srv := server.New(log, cfg)
	srv.CloseOnProgramEnd()
	srv.Start()
}
