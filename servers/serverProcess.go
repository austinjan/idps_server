package servers

import (
	"context"
	"github.com/austinjan/idps_server/servers/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func Run(ctx context.Context) {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://localhost:3001", "http://localhost:3000", "http://localhost:4000", "ws://localhost:3000", "ws://localhost:3001"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	settings := viper.GetStringMap("host")
	_addr, ok := settings["port"]
	if !ok {
		_addr = ":3011"
	}
	_timeout, ok := settings["timeout"]
	if !ok {
		_timeout = 5
	}
	srv := &http.Server{
		Addr:         _addr.(string),
		WriteTimeout: time.Second * _timeout.(time.Duration),
		ReadTimeout:  time.Second * _timeout.(time.Duration),
		Handler:      handler,
	}
	router.InitRouter(r)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Http server started..")

	// go mqttserver.Run()
	// fmt.Println("Running mqtt subscribe server...")

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	srv.Shutdown(ctx)
}
