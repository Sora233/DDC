package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Sora233/DDC/db"
	"github.com/Sora233/DDC/server/config"
	"github.com/Sora233/DDC/server/server"
	"github.com/sirupsen/logrus"
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/server_config.toml", "")

	flag.Parse()

	_, err := toml.DecodeFile(configPath, &config.Global)
	if err != nil {
		logrus.Fatalf("can not decode config: %v", err)
	}
	if lvl, err := logrus.ParseLevel(config.Global.LogLevel); err == nil {
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Warnf("unknown log level <%v>, use debug instead.", config.Global.LogLevel)
	}
}
func main() {
	if err := db.Init(&config.Global.DB); err != nil {
		logrus.Fatalf("init database error %v", err)
	}
	server.RunServer()
}
