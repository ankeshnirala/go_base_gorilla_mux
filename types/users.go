package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	State   string `json:"state" bson:"state,omitempty"`
	City    string `json:"city" bson:"city,omitempty"`
	Pincode uint32 `json:"pincode" bson:"pincode,omitempty"`
}

type User struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name,omitempty"`
	Age     int8               `json:"age" bson:"age,omitempty"`
	Address Address            `json:"address" bson:"address"`
}

type CreateUserResponse struct {
	InsertedID primitive.ObjectID `json:"insertID" bson:"insertID"`
}

type DeleteUserResponse struct {
	DeletedId   primitive.ObjectID `json:"deletedID" bson:"deletedID"`
	DeleteCount int8               `json:"deleteCount" bson:"deleteCount"`
}

type UpdateUserResponse struct {
	UpdateId    primitive.ObjectID `json:"UpdateID" bson:"UpdateID"`
	UpdateCount int8               `json:"UpdateCount" bson:"UpdateCount"`
}

func NewUser(name string, age int8, addresss Address) *User {
	return &User{
		Name:    name,
		Age:     age,
		Address: addresss,
	}
}
