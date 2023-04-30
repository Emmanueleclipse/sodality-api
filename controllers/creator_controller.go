package controllers

import (
	"context"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GetCreatorDirectoryByDirectoryName = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var allContent []*models.Content

	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "category_name", Value: params["category_name"]}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("contents does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.Content
		err := cursor.Decode(&content)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allContent = append(allContent, &content)
	}

	middlewares.SuccessArrRespond(allContent, rw)
})

var GetAllCreatorsContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var allContent []*models.GetAllContentWithCreatorResp
	// opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.GetAllContentWithCreatorResp
		err := cursor.Decode(&content)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allContent = append(allContent, &content)
	}

	for _, v := range allContent {
		var user models.User

		userID, _ := primitive.ObjectIDFromHex(v.UserID)
		collection := client.Database("sodality").Collection("users")
		collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
		user.Password = ""
		user.OTPEnabled = false
		user.OTPSecret = ""
		user.OTPAuthURL = ""
		v.User = user
	}
	middlewares.SuccessArrRespond(allContent, rw)
})
