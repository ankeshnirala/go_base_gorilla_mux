package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankeshnirala/go/go_services/constants"
	"github.com/ankeshnirala/go/go_services/middleware"
	"github.com/ankeshnirala/go/go_services/storage"
	"github.com/ankeshnirala/go/go_services/types"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage struct {
	mongoStore storage.MongoStore
	logger     *log.Logger
}

func New(logger *log.Logger, mongoStore storage.MongoStore) *Storage {
	return &Storage{
		mongoStore: mongoStore,
		logger:     logger,
	}
}

func (s *Storage) CreateUser(w http.ResponseWriter, r *http.Request) error {
	// sync request body data
	req := new(types.User)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.logger.Println(err.Error())
		return middleware.WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: constants.EMPTY_BODY})
	}

	user := types.NewUser(req.Name, req.Age, req.Address)
	result, err := s.mongoStore.InsertOne(constants.USER_COLLECTION, user)
	if err != nil {
		return err
	}

	resp := types.CreateUserResponse{
		InsertedID: result.InsertedID.(primitive.ObjectID),
	}

	return middleware.WriteJSON(w, http.StatusOK, resp)
}

func (s *Storage) GetUser(w http.ResponseWriter, r *http.Request) error {
	// name := r.URL.Query().Get("key")
	param := mux.Vars(r)["key"]

	// Decode the string ID from hexadecimal to bytes.
	objectId, _ := primitive.ObjectIDFromHex(param)
	query := bson.D{{Key: "_id", Value: objectId}}

	var user *types.User
	err := s.mongoStore.FindOne(constants.USER_COLLECTION, query).Decode(&user)

	if err != nil {
		s.logger.Println(err.Error())
		return middleware.WriteJSON(w, http.StatusNotFound, types.ApiError{Error: constants.NOT_FOUND})
	}

	return middleware.WriteJSON(w, http.StatusOK, user)
}

func (s *Storage) GetUsers(w http.ResponseWriter, r *http.Request) error {
	cursor, err := s.mongoStore.Find(constants.USER_COLLECTION, bson.M{})

	if err != nil {
		return middleware.WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: constants.WENT_WRONG})
	}

	var users []*types.User
	for cursor.Next(context.Background()) {
		var user *types.User
		err = cursor.Decode(&user)
		if err != nil {
			return middleware.WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: constants.WENT_WRONG})
		}

		users = append(users, user)
	}

	return middleware.WriteJSON(w, http.StatusOK, users)
}

func (s *Storage) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	// name := r.URL.Query().Get("key")
	param := mux.Vars(r)["key"]
	if param == "" {
		return middleware.WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: constants.EMPTY_PARAM})
	}

	// Decode the string ID from hexadecimal to bytes.
	objectId, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return middleware.WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: constants.EMPTY_PARAM})
	}

	req := new(types.User)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.logger.Println(err.Error())
		return middleware.WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: constants.EMPTY_BODY})
	}

	// Create a filter to find the document to update.
	filter := bson.M{"_id": objectId}

	var upd bson.M
	pByte, _ := bson.Marshal(req)
	bson.Unmarshal(pByte, &upd)
	update := bson.D{{Key: "$set", Value: upd}}

	result, err := s.mongoStore.UpdateOne(constants.USER_COLLECTION, filter, update)
	if err != nil {
		s.logger.Println(err.Error())
		return middleware.WriteJSON(w, http.StatusNotFound, types.ApiError{Error: constants.UPDATEFAILED})
	}

	resp := types.UpdateUserResponse{
		UpdateId:    objectId,
		UpdateCount: int8(result.ModifiedCount),
	}

	return middleware.WriteJSON(w, http.StatusOK, resp)
}

func (s *Storage) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	// name := r.URL.Query().Get("key")
	param := mux.Vars(r)["key"]

	// Decode the string ID from hexadecimal to bytes.
	objectId, _ := primitive.ObjectIDFromHex(param)
	query := bson.D{{Key: "_id", Value: objectId}}

	result, err := s.mongoStore.DeleteByID(constants.USER_COLLECTION, query)
	if err != nil {
		s.logger.Println(err.Error())
		return middleware.WriteJSON(w, http.StatusNotFound, types.ApiError{Error: constants.NOT_FOUND})
	}

	resp := types.DeleteUserResponse{
		DeletedId:   objectId,
		DeleteCount: int8(result.DeletedCount),
	}

	return middleware.WriteJSON(w, http.StatusOK, resp)
}
