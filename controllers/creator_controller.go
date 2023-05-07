package controllers

import (
	"context"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"
	"time"

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

var GetAllCreators = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var allCreators []*models.GetAllCreatorsResp
	// opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "total_donations", Value: -1}})

	collection := client.Database("sodality").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "role", Value: "creator"}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.SuccessArrRespond(nil, rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var creator models.GetAllCreatorsResp
		err := cursor.Decode(&creator)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		pipeline := bson.A{
			bson.M{"$match": bson.M{"creator_id": creator.ID.Hex(), "expired_at": bson.M{"$gte": time.Now().UTC()}}},
			bson.M{"$group": bson.M{"_id": "$user_id", "count": bson.M{"$sum": 1}}},
		}

		var supporterCount []bson.M
		collection := client.Database("sodality").Collection("donations")
		cur, err := collection.Aggregate(context.TODO(), pipeline)
		if err != nil && err != mongo.ErrNoDocuments {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}
			supporterCount = append(supporterCount, result)
		}

		creator.Supporters = int64(len(supporterCount))

		allCreators = append(allCreators, &creator)
	}
	middlewares.SuccessArrRespond(allCreators, rw)
})

var SearchCreatorByUsername = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var allCreator []*models.GetAllCreatorsResp
	if len(params["search"]) < 2 {
		middlewares.SuccessArrRespond(nil, rw)
		// middlewares.ErrorResponse("search required two or more alphabets or numbers", rw)
		return
	}

	filter := bson.M{"username": bson.M{"$regex": "^" + params["search"], "$options": "im"}}

	collection := client.Database("sodality").Collection("users")
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
		var creator models.GetAllCreatorsResp
		err := cursor.Decode(&creator)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		pipeline := bson.A{
			bson.M{"$match": bson.M{"creator_id": creator.ID.Hex(), "expired_at": bson.M{"$gte": time.Now().UTC()}}},
			bson.M{"$group": bson.M{"_id": "$user_id", "count": bson.M{"$sum": 1}}},
		}

		var supporterCount []bson.M
		collection := client.Database("sodality").Collection("donations")
		cur, err := collection.Aggregate(context.TODO(), pipeline)
		if err != nil && err != mongo.ErrNoDocuments {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}
			supporterCount = append(supporterCount, result)
		}

		creator.Supporters = int64(len(supporterCount))

		allCreator = append(allCreator, &creator)
	}

	middlewares.SuccessArrRespond(allCreator, rw)
})
