package routers

import (
	"log"

	"github.com/ankeshnirala/go/go_services/constants"
	"github.com/ankeshnirala/go/go_services/handlers"
	"github.com/ankeshnirala/go/go_services/middleware"
	"github.com/ankeshnirala/go/go_services/storage"
	"github.com/gorilla/mux"
)

type Storage struct {
	logger     *log.Logger
	mongoStore storage.MongoStore
	userRouter *mux.Router
}

func New(logger *log.Logger, mongoStore storage.MongoStore, userRouter *mux.Router) *Storage {
	return &Storage{
		mongoStore: mongoStore,
		userRouter: userRouter,
		logger:     logger,
	}
}

func (s *Storage) RegisterUserRoutes() {
	userController := handlers.New(s.logger, s.mongoStore)

	s.userRouter.HandleFunc(constants.ADD, middleware.MakeHTTPHandleFunc(userController.CreateUser)).Methods("POST")
	s.userRouter.HandleFunc(constants.GET, middleware.MakeHTTPHandleFunc(userController.GetUser)).Methods("GET")
	s.userRouter.HandleFunc(constants.GETALL, middleware.MakeHTTPHandleFunc(userController.GetUsers)).Methods("GET")
	s.userRouter.HandleFunc(constants.UPDATE, middleware.MakeHTTPHandleFunc(userController.UpdateUser)).Methods("PUT")
	s.userRouter.HandleFunc(constants.DELETE, middleware.MakeHTTPHandleFunc(userController.DeleteUser)).Methods("DELETE")
}
