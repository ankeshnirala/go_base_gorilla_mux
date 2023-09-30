package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankeshnirala/go/go_services/routers"
	"github.com/ankeshnirala/go/go_services/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	logger     *log.Logger
	listenAddr string
	mongoStore storage.MongoStore
}

func NewServer(logger *log.Logger, listenAddr string, mongoStore storage.MongoStore) *Server {
	return &Server{
		logger:     logger,
		listenAddr: listenAddr,
		mongoStore: mongoStore,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	// Create a subrouter for the user router.
	userRouter := router.PathPrefix("/users").Subrouter()
	routers.New(s.logger, s.mongoStore, userRouter).RegisterUserRoutes()

	router.MethodNotAllowedHandler = MakeHTTPHandleFunc(s.MethodNotAllowed)
	router.NotFoundHandler = MakeHTTPHandleFunc(s.PathNotFound)

	srv := &http.Server{
		Addr:         s.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Println(err, "server.go")
		}
	}()

	gracefulShutdown(s.logger, srv)

	return nil
}

func gracefulShutdown(l *log.Logger, s *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown...", sig)

	tc, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		l.Fatal(err)
	}
	s.Shutdown(tc)
}
