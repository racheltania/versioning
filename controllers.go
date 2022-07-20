package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var userCollection = db().Database("cobadb").Collection("people")

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	result, _ := userCollection.InsertOne(context.TODO(), person)
	res := Result{Code: 200, Data: result, Message: "success create data"}
	hasil, err := json.Marshal(res)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(hasil)
}

func DeletePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	result, err := userCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}
	res := Result{Code: 200, Data: result, Message: "success delete data"}
	hasil, err := json.Marshal(res)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(hasil)

}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	err := userCollection.FindOne(context.TODO(), Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	res := Result{Code: 200, Data: person, Message: "success get data"}
	hasil, err := json.Marshal(res)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(hasil)
}

func UpdatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var updateperson Person
	_ = json.NewDecoder(request.Body).Decode(&updateperson)

	result, err := userCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.D{
		{"$set", &updateperson}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	res := Result{Code: 200, Data: result, Message: "success update data"}
	hasil, err := json.Marshal(res)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(hasil)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	var people []Person
	cursor, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	res := Result{Code: 200, Data: people, Message: "success get data"}
	hasil, err := json.Marshal(res)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(hasil)
}
