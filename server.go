package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/austinjan/idps_server/config"
	"github.com/austinjan/idps_server/servers"
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

	httpCtx, httpDone := context.WithCancel(context.Background())
	go servers.Run(httpCtx)

	// db := mongodb.GetDB()
	// db.SaveTagPosition(bson.M{"test": "testtext"})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	httpDone()
	os.Exit(0)
}
