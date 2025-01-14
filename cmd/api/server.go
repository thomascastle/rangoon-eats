package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Info("server started", map[string]string{"addr": server.Addr, "env": app.config.env})

	return server.ListenAndServe()
}
