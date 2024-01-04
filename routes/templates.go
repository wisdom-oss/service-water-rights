package routes

import (
	"errors"
	"net/http"

	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v3"
)

// BasicHandler contains just a response, that is used to show the templating
func BasicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there"))
}

func BasicWithErrorHandling(w http.ResponseWriter, r *http.Request) {
	// access the error handlers
	errorHandler := r.Context().Value(wisdomMiddleware.ERROR_CHANNEL_NAME).(chan interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.STATUS_CHANNEL_NAME).(chan bool)
	// now publish an error to each of the wisdom errors
	errorHandler <- errors.New("native test error")
	// now wait for the error to be handled
	<-statusChannel
	return
}
