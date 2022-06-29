package main

import (
	"fmt"
	"go-redis/config"
	"go-redis/tcp"
	"go.uber.org/zap"
	"log"
	"os"
)

const configFile = "redis.conf"

var defaultConfig = &config.ServerConfig{Addr: "0.0.0.0", Port: 6379}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err != nil && !info.IsDir()
}

func main() {
	// init logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot init zap logger")
	}
	zap.ReplaceGlobals(logger)

	// init config
	if fileExists(configFile) {
		config.SetConfig(configFile)
	} else {
		config.Config = defaultConfig
	}

	// init tcp server
	tcpServer := tcp.NewServer(fmt.Sprintf("%s:%d", config.Config.Addr, config.Config.Port))
	if err = tcpServer.ListenAndServeWithSignal(tcp.NewTestHandler()); err != nil {
		zap.S().Fatalf("cannot init tcp server")
	}
}
