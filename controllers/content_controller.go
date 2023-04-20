package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostContent -> Create a creator content
var PostContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var content models.Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var existingUser models.User
	userID, _ := primitive.ObjectIDFromHex(props["user_id"].(string))

	userCollection := client.Database("sodality").Collection("users")
	err = userCollection.FindOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if existingUser.ID != userID || err == mongo.ErrNoDocuments {
		middlewares.ErrorResponse("user does not exists", rw)
		return
	}

	content.UserID = userID.Hex()
	content.CreatedAt = time.Now().UTC()

	contentCollection := client.Database("sodality").Collection("content")
	result, err := contentCollection.InsertOne(context.TODO(), content)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

// GetContentByID -> Get content of user by content id
var GetContentByID = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var content models.Content

	contentID, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("sodality").Collection("content")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: contentID}}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content id does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessArrRespond(content, rw)
})

var GetOwnContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var allContent []*models.Content

	// opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: props["user_id"].(string)}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.Content
		err := cursor.Decode(&content)
		if err != nil {
			log.Fatal(err)
		}

		allContent = append(allContent, &content)
	}

	middlewares.SuccessArrRespond(allContent, rw)
})

var SearchContentByTitle = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var allContent []*models.Content

	filter := bson.M{"title": bson.M{"$regex": params["search"], "$options": "im"}}

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), filter)
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
			log.Fatal(err)
		}

		allContent = append(allContent, &content)
	}

	middlewares.SuccessArrRespond(allContent, rw)
})
