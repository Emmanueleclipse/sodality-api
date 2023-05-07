package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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
	content.Locked = true
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
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allContent = append(allContent, &content)
	}

	middlewares.SuccessArrRespond(allContent, rw)
})

var GetCreatorContentById = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var allContent []*models.GetContentResp

	// opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: params["creator_id"]}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.GetContentResp
		err := cursor.Decode(&content)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allContent = append(allContent, &content)
	}

	middlewares.SuccessArrRespond(allContent, rw)
})

var SearchContentByTitle = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	var allContent []*models.Content

	search := r.URL.Query().Get("search")
	limit := r.URL.Query().Get("limit")

	if len(search) < 2 {
		middlewares.SuccessArrRespond(nil, rw)
		// middlewares.ErrorResponse("search required two or more alphabets or numbers", rw)
		return
	}
	var newLimit int64
	var err error
	if len(limit) > 0 {
		newLimit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

	}
	opts := options.Find().SetLimit(newLimit)

	filter := bson.M{"title": bson.M{"$regex": "^" + search, "$options": "im"}}

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), filter, opts)
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

var DeleteContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	params := mux.Vars(r)

	contentID, _ := primitive.ObjectIDFromHex(params["id"])
	userID := props["user_id"].(string)

	filter := bson.M{"$and": []interface{}{
		bson.M{"_id": contentID},
		bson.M{"user_id": userID}}}

	collection := client.Database("sodality").Collection("content")
	del, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	if del.DeletedCount == 0 {
		middlewares.ErrorResponse("you have no access to delete this content", rw)
		return
	}

	middlewares.SuccessResponse("content delete successfully", rw)
})
