package main

import (
	"context"
	"fmt"
	_ "github.com/austinjan/idps_server/config"
	"github.com/austinjan/idps_server/servers"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
)

func main() {
	logFile, err := os.Create("log")
	fmt.Println("idps server start version: ", viper.GetString("version"))
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	httpCtx, httpDone := context.WithCancel(context.Background())
	go servers.Run(httpCtx)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	httpDone()
	os.Exit(0)
}
