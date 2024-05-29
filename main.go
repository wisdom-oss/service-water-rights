package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	errorHandler "github.com/wisdom-oss/microservice-middlewares/v5/error"
	securityMiddlewares "github.com/wisdom-oss/microservice-middlewares/v5/security"
	"golang.org/x/net/http2"

	"github.com/wisdom-oss/service-water-rights/globals"
	"github.com/wisdom-oss/service-water-rights/routes"

	"golang.org/x/net/http2/h2c"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	mainLogger := log.With().Str("service", globals.ServiceName).Logger()
	mainLogger.Info().Msg("bootstrapping http server")

	// create a new chi.Router which handles the routing to the different routes
	router := chi.NewRouter()
	// add a middleware that uses the x-real-ip or x-forwarded-for headers to
	// show the real ip of the person sending a request
	router.Use(config.Middlewares...)
	// now configure the logging for the service and add it to the router
	httplog.Configure(httplog.Options{
		JSON:    true,
		Concise: true,
	})

	// now add the routes and their path specifications to the router
	router.
		With(securityMiddlewares.RequireScope(globals.ServiceName, securityMiddlewares.ScopeRead)).
		Get("/", routes.UsageLocations)
	//router.Get("/details/{water-right-nlwkn-id}", routes.SingleWaterRight)

	// now configure the http2c and the http server
	http2Server := &http2.Server{}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.ListenAddress, globals.Environment["LISTEN_PORT"]),
		Handler: h2c.NewHandler(router, http2Server),
	}

	// now setup some signal handling to allow stopping the service gracefully
	cancelQueue := make(chan os.Signal, 1)
	signal.Notify(cancelQueue, os.Interrupt)

	// now start up the http server
	go func() {
		mainLogger.Info().Msg("starting http server")
		if err := httpServer.ListenAndServe(); err != nil {
			mainLogger.Fatal().Err(err).Msg("unable to start http server")
		}
	}()

	// wait for the cancel signal
	<-cancelQueue
	// shutdown the http server gracefully
	if err := httpServer.Shutdown(nil); err != nil {
		mainLogger.Fatal().Err(err).Msg("unable to shutdown http server")
	}
}
