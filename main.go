package main

import (
	"fmt"
	"go_redis/config"
	"go_redis/lib/logger"
	"go_redis/tcp"
	"os"
)

const configPath string = "redis.conf"

var defaultConfig = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configPath) {
		config.SetupConfig(configPath)
	} else {
		config.Properties = defaultConfig
	}
	handler := tcp.NewHandler()
	err := tcp.ListenAndServerWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, handler)
	if err != nil {
		logger.Error(err)
	}
}
