package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	_ "github.com/austinjan/config"
	"github.com/spf13/viper"
)

func main() {
	logFile, err := os.Create("log")
	fmt.Println("idps server start version: ", viper.GetString("version"))
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	os.Exit(0)
}